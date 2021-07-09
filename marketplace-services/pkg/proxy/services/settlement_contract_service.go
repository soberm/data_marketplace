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

type SettlementContractService interface {
	Deposit(ctx context.Context, amount *big.Int) (*types.Transaction, error)
	SettleTrade(ctx context.Context, counter *big.Int) (*types.Transaction, error)
	ResolveDispute(ctx context.Context, counter *big.Int) (*types.Transaction, error)
	ResolveTimeout(ctx context.Context) (*types.Transaction, error)
	GetProviderCounter(ctx context.Context) (*contracts.Counter, error)
	GetConsumerCounter(ctx context.Context) (*contracts.Counter, error)
	GetBrokerCounter(ctx context.Context) (*contracts.Counter, error)
	GetSettlement(ctx context.Context) (*contracts.Settlement, error)
}

type settlementContractServiceImpl struct {
	logger             logrus.FieldLogger
	keyStore           *keystore.KeyStore
	walletService      WalletService
	settlementContract contracts.SettlementContract
}

func NewSettlementContractServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	settlementContract contracts.SettlementContract,
) *settlementContractServiceImpl {
	return &settlementContractServiceImpl{
		logger:             logger,
		walletService:      walletService,
		keyStore:           keyStore,
		settlementContract: settlementContract,
	}
}

func (s settlementContractServiceImpl) Deposit(ctx context.Context, amount *big.Int) (*types.Transaction, error) {
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
	transactOpts.Value = amount

	if err := s.keyStore.Unlock(account, w.Passphrase); err != nil {
		return nil, fmt.Errorf("unlock account %s: %w", account.Address.Hex(), err)
	}
	defer func() {
		if err := s.keyStore.Lock(account.Address); err != nil {
			s.logger.Warnf("lock account %s: %w", account.Address.Hex(), err)
		}
	}()

	return s.settlementContract.Deposit(transactOpts)
}

func (s settlementContractServiceImpl) SettleTrade(ctx context.Context, counter *big.Int) (*types.Transaction, error) {
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

	return s.settlementContract.SettleTrade(transactOpts, counter)
}

func (s settlementContractServiceImpl) ResolveDispute(ctx context.Context, counter *big.Int) (*types.Transaction, error) {
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

	return s.settlementContract.ResolveDispute(transactOpts, counter)
}

func (s settlementContractServiceImpl) ResolveTimeout(ctx context.Context) (*types.Transaction, error) {
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

	return s.settlementContract.ResolveTimeout(transactOpts)
}

func (s settlementContractServiceImpl) GetProviderCounter(ctx context.Context) (*contracts.Counter, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.settlementContract.GetProviderCounter(callOpts)
}

func (s settlementContractServiceImpl) GetConsumerCounter(ctx context.Context) (*contracts.Counter, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.settlementContract.GetConsumerCounter(callOpts)
}

func (s settlementContractServiceImpl) GetBrokerCounter(ctx context.Context) (*contracts.Counter, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.settlementContract.GetBrokerCounter(callOpts)
}

func (s settlementContractServiceImpl) GetSettlement(ctx context.Context) (*contracts.Settlement, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.settlementContract.GetSettlement(callOpts)
}
