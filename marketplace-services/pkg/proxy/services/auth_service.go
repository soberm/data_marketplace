package services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"marketplace-services/pkg/proxy/model"
	"time"
)

type AuthService interface {
	GetToken(ctx context.Context, username string, password []byte) (string, error)
	GenerateToken(user *model.Account) (string, error)
	ParseToken(token string) (*CustomClaims, error)
	AuthFunction() func(ctx context.Context) (context.Context, error)
}

type CustomClaims struct {
	jwt.StandardClaims
	UserID   uint       `json:"uid,omitempty"`
	UserRole model.Role `json:"role,omitempty"`
}

type authServiceImpl struct {
	logger         logrus.FieldLogger
	userService    AccountService
	appName        string
	signingKey     []byte
	expirationTime int64
}

func NewAuthServiceImpl(
	logger logrus.FieldLogger,
	userService AccountService,
	appName string,
	signingKey []byte,
	expirationTime int64,
) *authServiceImpl {
	return &authServiceImpl{
		logger:         logger,
		userService:    userService,
		appName:        appName,
		signingKey:     signingKey,
		expirationTime: expirationTime,
	}
}

func (s *authServiceImpl) GetToken(ctx context.Context, username string, password []byte) (string, error) {
	u, err := s.userService.FindAccountByNameAndPassword(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}
	token, err := s.GenerateToken(u)
	if err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}
	return token, nil
}

func (s *authServiceImpl) GenerateToken(user *model.Account) (string, error) {
	s.logger.Infof("Generating token for user %s", user.Name)

	now := time.Now()
	customClaims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: now.Unix() + s.expirationTime,
			Id:        uuid.New(),
			IssuedAt:  now.Unix(),
			Issuer:    s.appName,
			NotBefore: now.Unix(),
			Subject:   user.Name,
		},
		UserID:   user.ID,
		UserRole: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", fmt.Errorf("get token for user %s: %w", user.Name, err)
	}
	return tokenString, err

}
func (s *authServiceImpl) ParseToken(tokenString string) (*CustomClaims, error) {
	s.logger.Debugf("Parsing token %s", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.signingKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("parse token: %w", err)
	}
}

func (s *authServiceImpl) AuthFunction() func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		claims, err := s.ParseToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failure: %w", err)
		}

		u := model.Account{Model: gorm.Model{
			ID: claims.UserID,
		},
			Name: claims.Subject,
			Role: claims.UserRole,
		}
		return context.WithValue(ctx, "principal", u), nil
	}
}
