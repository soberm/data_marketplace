package services

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"marketplace-services/pkg/proxy/model"
)

type AccountService interface {
	CreateAccount(ctx context.Context, acc *model.Account) (*model.Account, error)
	FindAccountByName(ctx context.Context, name string) (*model.Account, error)
	FindAccountByNameAndPassword(ctx context.Context, name string, password []byte) (*model.Account, error)
	FindAccountById(ctx context.Context, id uint) (*model.Account, error)
	FindAccounts(ctx context.Context) ([]*model.Account, error)
}

type accountServiceImpl struct {
	db     *gorm.DB
	logger logrus.FieldLogger
}

func NewAccountServiceImpl(db *gorm.DB, logger logrus.FieldLogger) *accountServiceImpl {
	return &accountServiceImpl{
		db:     db,
		logger: logger,
	}
}

func (s *accountServiceImpl) CreateAccount(ctx context.Context, acc *model.Account) (*model.Account, error) {
	account, ok := ctx.Value("principal").(model.Account)
	if !ok {
		return nil, fmt.Errorf("extract account from context")
	}

	if !account.HasRole(model.RoleAdmin) {
		return nil, fmt.Errorf("role %d needed", model.RoleAdmin)
	}

	err := acc.Validate()
	if err != nil {
		return nil, fmt.Errorf("validate account: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(acc.Password, bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("generate hash from password: %w", err)
	}
	acc.Password = hashedPassword

	return acc, s.db.Create(acc).Error
}

func (s *accountServiceImpl) FindAccountByName(_ context.Context, name string) (*model.Account, error) {
	var user model.Account
	err := s.db.Preload("Wallet").Where(&model.Account{Name: name}).First(&user).Error
	if err != nil {
		err = fmt.Errorf("preload wallet and get first account with name %s: %w", name, err)
	}
	return &user, err
}

func (s *accountServiceImpl) FindAccountByNameAndPassword(ctx context.Context, name string, password []byte) (*model.Account, error) {
	account, err := s.FindAccountByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("find account by name %s: %w", name, err)
	}

	err = bcrypt.CompareHashAndPassword(account.Password, password)
	if err != nil {
		return nil, fmt.Errorf("compare hash and password of account %s: %w", name, err)
	}
	return account, err
}

func (s *accountServiceImpl) FindAccountById(_ context.Context, id uint) (*model.Account, error) {
	var account model.Account
	err := s.db.Preload("Wallet").First(&account, id).Error
	if err != nil {
		err = fmt.Errorf("preload wallet and get first account with id %d: %w", id, err)
	}
	return &account, err
}

func (s *accountServiceImpl) FindAccounts(_ context.Context) ([]*model.Account, error) {
	var accounts []*model.Account
	err := s.db.Preload("Wallet").Find(&accounts).Error
	if err != nil {
		err = fmt.Errorf("preload wallet and get all accounts: %w", err)
	}
	return accounts, err
}
