package api

import (
	"context"
	"marketplace-services/pkg/proxy/services"
)

type authServiceServer struct {
	UnimplementedAuthServiceServer
	authService services.AuthService
}

func NewAuthServiceServer(authService services.AuthService) *authServiceServer {
	return &authServiceServer{authService: authService}
}

func (s *authServiceServer) GetToken(ctx context.Context, req *GetTokenRequest) (*GetTokenResponse, error) {
	token, err := s.authService.GetToken(ctx, req.Username, req.Password)
	return &GetTokenResponse{Token: token}, err
}

func (s *authServiceServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if fullMethodName == "/proxy.AuthService/GetToken" {
		return ctx, nil
	}
	return s.authService.AuthFunction()(ctx)
}
