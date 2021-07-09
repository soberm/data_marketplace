package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"marketplace-services/pkg/contracts"
	"math/big"
)

type UserContractService interface {
	CreateUser(ctx context.Context, user *contracts.User) (*types.Transaction, error)
	UpdateUser(ctx context.Context, user *contracts.User) (*types.Transaction, error)
	RemoveUser(ctx context.Context) (*types.Transaction, error)
	FindUserByIndex(ctx context.Context, index *big.Int) (*contracts.User, error)
	FindUserByAddress(ctx context.Context, address common.Address) (*contracts.User, error)
	CountUsers(ctx context.Context) (*big.Int, error)
	ExistsUserByAddress(ctx context.Context, address common.Address) (bool, error)
	ExistsUserByAddressAndDeleted(ctx context.Context, address common.Address, deleted bool) (bool, error)
}

type userContractServiceImpl struct {
	logger        logrus.FieldLogger
	keyStore      *keystore.KeyStore
	walletService WalletService
	userContract  contracts.UserContract
}

func NewUserContractServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	userContract contracts.UserContract,
) *userContractServiceImpl {
	return &userContractServiceImpl{
		logger:        logger,
		walletService: walletService,
		keyStore:      keyStore,
		userContract:  userContract,
	}
}

func (s *userContractServiceImpl) CreateUser(ctx context.Context, user *contracts.User) (*types.Transaction, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}

	account := accounts.Account{Address: common.BytesToAddress(w.Address)}
	transactOpts, err := bind.NewKeyStoreTransactor(
		s.keyStore,
		account,
	)
	if err != nil {
		return nil, fmt.Errorf("new keystore transactor: %w", err)
	}

	if err := s.keyStore.Unlock(account, w.Passphrase); err != nil {
		return nil, fmt.Errorf("unlock account %s: %w", account.Address.Hex(), err)
	}
	defer func() {
		if err := s.keyStore.Lock(account.Address); err != nil {
			s.logger.Warnf("lock account %s: %w", account.Address.Hex(), err)
		}
	}()

	return s.userContract.CreateUser(transactOpts, user)
}

func (s *userContractServiceImpl) UpdateUser(ctx context.Context, user *contracts.User) (*types.Transaction, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}

	account := accounts.Account{Address: common.BytesToAddress(w.Address)}
	transactOpts, err := bind.NewKeyStoreTransactor(
		s.keyStore,
		account,
	)
	if err != nil {
		return nil, fmt.Errorf("new keystore transactor: %w", err)
	}

	if err := s.keyStore.Unlock(account, w.Passphrase); err != nil {
		return nil, fmt.Errorf("unlock account %s: %w", account.Address.Hex(), err)
	}
	defer func() {
		if err := s.keyStore.Lock(account.Address); err != nil {
			s.logger.Warnf("lock account %s: %w", account.Address.Hex(), err)
		}
	}()

	return s.userContract.UpdateUser(transactOpts, user)
}

func (s *userContractServiceImpl) RemoveUser(ctx context.Context) (*types.Transaction, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}

	account := accounts.Account{Address: common.BytesToAddress(w.Address)}
	transactOpts, err := bind.NewKeyStoreTransactor(
		s.keyStore,
		account,
	)
	if err != nil {
		return nil, fmt.Errorf("new keystore transactor: %w", err)
	}

	if err := s.keyStore.Unlock(account, w.Passphrase); err != nil {
		return nil, fmt.Errorf("unlock account %s: %w", account.Address.Hex(), err)
	}
	defer func() {
		if err := s.keyStore.Lock(account.Address); err != nil {
			s.logger.Warnf("lock account %s: %w", account.Address.Hex(), err)
		}
	}()

	return s.userContract.RemoveUser(transactOpts)
}

func (s *userContractServiceImpl) FindUserByIndex(ctx context.Context, index *big.Int) (*contracts.User, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.userContract.FindUserByIndex(callOpts, index)
}

func (s *userContractServiceImpl) FindUserByAddress(
	ctx context.Context,
	address common.Address,
) (*contracts.User, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.userContract.FindUserByAddress(callOpts, address)
}

func (s *userContractServiceImpl) CountUsers(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.userContract.CountUsers(callOpts)
}

func (s *userContractServiceImpl) ExistsUserByAddress(ctx context.Context, address common.Address) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.userContract.ExistsUserByAddress(callOpts, address)
}

func (s *userContractServiceImpl) ExistsUserByAddressAndDeleted(
	ctx context.Context,
	address common.Address,
	deleted bool,
) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.userContract.ExistsUserByAddressAndDeleted(callOpts, address, deleted)
}
