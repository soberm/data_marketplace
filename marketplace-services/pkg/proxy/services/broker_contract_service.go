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

type BrokerContractService interface {
	CreateBroker(ctx context.Context, broker *contracts.Broker) (*types.Transaction, error)
	UpdateBroker(ctx context.Context, broker *contracts.Broker) (*types.Transaction, error)
	RemoveBroker(ctx context.Context, address common.Address) (*types.Transaction, error)
	FindBrokerByIndex(ctx context.Context, index *big.Int) (*contracts.Broker, error)
	FindBrokerByAddress(ctx context.Context, address common.Address) (*contracts.Broker, error)
	CountBrokers(ctx context.Context) (*big.Int, error)
	ExistsBrokerByAddress(ctx context.Context, address common.Address) (bool, error)
	ExistsBrokerByAddressAndDeleted(ctx context.Context, address common.Address, deleted bool) (bool, error)
}

type brokerContractServiceImpl struct {
	logger         logrus.FieldLogger
	keyStore       *keystore.KeyStore
	walletService  WalletService
	brokerContract contracts.BrokerContract
}

func NewBrokerContractServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	brokerContract contracts.BrokerContract,
) *brokerContractServiceImpl {
	return &brokerContractServiceImpl{
		logger:         logger,
		walletService:  walletService,
		keyStore:       keyStore,
		brokerContract: brokerContract,
	}
}

func (s brokerContractServiceImpl) CreateBroker(
	ctx context.Context,
	broker *contracts.Broker,
) (*types.Transaction, error) {
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

	return s.brokerContract.CreateBroker(transactOpts, broker)
}

func (s brokerContractServiceImpl) UpdateBroker(
	ctx context.Context,
	broker *contracts.Broker,
) (*types.Transaction, error) {
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

	return s.brokerContract.UpdateBroker(transactOpts, broker)
}

func (s brokerContractServiceImpl) RemoveBroker(
	ctx context.Context,
	address common.Address,
) (*types.Transaction, error) {
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

	return s.brokerContract.RemoveBroker(transactOpts, address)
}

func (s brokerContractServiceImpl) FindBrokerByIndex(ctx context.Context, index *big.Int) (*contracts.Broker, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.brokerContract.FindBrokerByIndex(callOpts, index)
}

func (s brokerContractServiceImpl) FindBrokerByAddress(
	ctx context.Context,
	address common.Address,
) (*contracts.Broker, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.brokerContract.FindBrokerByAddress(callOpts, address)
}

func (s brokerContractServiceImpl) CountBrokers(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.brokerContract.CountBrokers(callOpts)
}

func (s brokerContractServiceImpl) ExistsBrokerByAddress(ctx context.Context, address common.Address) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.brokerContract.ExistsBrokerByAddress(callOpts, address)
}

func (s brokerContractServiceImpl) ExistsBrokerByAddressAndDeleted(
	ctx context.Context,
	address common.Address,
	deleted bool,
) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.brokerContract.ExistsBrokerByAddressAndDeleted(callOpts, address, deleted)
}
