package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"marketplace-services/pkg/proxy/model"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)
	FindWalletById(ctx context.Context, id uint) (*model.Wallet, error)
	FindWalletByAccountId(ctx context.Context, id uint) (*model.Wallet, error)
	FindWalletByAuthenticatedAccount(ctx context.Context) (*model.Wallet, error)
	FindKeyByAuthenticatedAccount(ctx context.Context) (*keystore.Key, error)
}

type walletServiceImpl struct {
	db       *gorm.DB
	logger   logrus.FieldLogger
	keyStore *keystore.KeyStore
}

func NewWalletServiceImpl(db *gorm.DB, logger logrus.FieldLogger, keyStore *keystore.KeyStore) *walletServiceImpl {
	return &walletServiceImpl{
		db:       db,
		logger:   logger,
		keyStore: keyStore,
	}
}

func (s *walletServiceImpl) CreateWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {

	proxyAccount, ok := ctx.Value("principal").(model.Account)
	if !ok {
		return nil, fmt.Errorf("extract account from context")
	}

	if !proxyAccount.HasRole(model.RoleAdmin) {
		return nil, fmt.Errorf("role %d needed", model.RoleAdmin)
	}

	err := wallet.Validate()
	if err != nil {
		return nil, fmt.Errorf("validate wallet: %w", err)
	}

	account, err := s.keyStore.NewAccount(wallet.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("create wallet for account %d: %w", wallet.AccountID, err)
	}

	wallet.Address = account.Address.Bytes()
	wallet.FilePath = account.URL.Path

	keyFile, err := ioutil.ReadFile(wallet.FilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file %s: %w", wallet.FilePath, err)
	}

	key, err := keystore.DecryptKey(keyFile, wallet.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt key of wallet %d: %w", wallet.ID, err)
	}

	wallet.PublicKey = crypto.FromECDSAPub(&key.PrivateKey.PublicKey)

	err = s.db.Create(wallet).Error
	if err != nil {
		if err := s.keyStore.Delete(account, wallet.Passphrase); err != nil {
			s.logger.Warnf("delete ethereum account %s", account.Address.Hex())
		}
		return nil, err
	}
	return wallet, err
}

func (s *walletServiceImpl) FindWalletById(ctx context.Context, id uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := s.db.First(&wallet, id).Error
	if err != nil {
		return nil, fmt.Errorf("get first wallet with id %d: %w", id, err)
	}
	return &wallet, err
}

func (s *walletServiceImpl) FindWalletByAccountId(ctx context.Context, id uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := s.db.Where(&model.Wallet{AccountID: id}).First(&wallet).Error
	if err != nil {
		return nil, fmt.Errorf("get first wallet of account %d: %w", id, err)
	}
	return &wallet, err
}

func (s *walletServiceImpl) FindKeyByAuthenticatedAccount(ctx context.Context) (*keystore.Key, error) {
	principal, ok := ctx.Value("principal").(model.Account)
	if !ok {
		return nil, fmt.Errorf("extract account from context")
	}

	var wallet model.Wallet
	err := s.db.Where(&model.Wallet{AccountID: principal.ID}).First(&wallet).Error
	if err != nil {
		return nil, fmt.Errorf("get first wallet of account %d: %w", principal.ID, err)
	}

	keyFile, err := ioutil.ReadFile(wallet.FilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file %s: %w", wallet.FilePath, err)
	}

	key, err := keystore.DecryptKey(keyFile, wallet.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt key of wallet %d: %w", wallet.ID, err)
	}

	return key, err
}

func (s *walletServiceImpl) FindWalletByAuthenticatedAccount(ctx context.Context) (*model.Wallet, error) {
	principal, ok := ctx.Value("principal").(model.Account)
	if !ok {
		return nil, fmt.Errorf("extract account from context")
	}

	wallet, err := s.FindWalletByAccountId(ctx, principal.ID)
	if err != nil {
		return nil, fmt.Errorf("find account by id %d: %w", principal.ID, err)
	}
	return wallet, nil
}
