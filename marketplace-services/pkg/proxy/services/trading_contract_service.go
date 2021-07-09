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

type TradingContractService interface {
	RequestTrading(ctx context.Context, product *big.Int, broker common.Address, startTime *big.Int, endTime *big.Int) (*types.Transaction, error)
	AcceptTradingRequest(ctx context.Context, id *big.Int) (*types.Transaction, error)
	DeclineTradingRequest(ctx context.Context, id *big.Int) (*types.Transaction, error)
	CreateTrade(ctx context.Context, negotiation *big.Int, address common.Address) (*types.Transaction, error)
	FindTradingRequestByIndex(ctx context.Context, index *big.Int) (*contracts.TradingRequest, error)
	FindTradeByIndex(ctx context.Context, index *big.Int) (*contracts.Trade, error)
	FindTradingRequestById(ctx context.Context, id *big.Int) (*contracts.TradingRequest, error)
	FindTradeById(ctx context.Context, id *big.Int) (*contracts.Trade, error)
	CountTradingRequests(ctx context.Context) (*big.Int, error)
	CountTrades(ctx context.Context) (*big.Int, error)
}

type tradingContractServiceImpl struct {
	logger          logrus.FieldLogger
	keyStore        *keystore.KeyStore
	walletService   WalletService
	tradingContract contracts.TradingContract
}

func NewTradingContractServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	tradingContract contracts.TradingContract,
) *tradingContractServiceImpl {
	return &tradingContractServiceImpl{
		logger:          logger,
		walletService:   walletService,
		keyStore:        keyStore,
		tradingContract: tradingContract,
	}
}

func (s tradingContractServiceImpl) RequestTrading(ctx context.Context, product *big.Int, broker common.Address, startTime *big.Int, endTime *big.Int) (*types.Transaction, error) {
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

	return s.tradingContract.RequestTrading(
		transactOpts,
		product,
		broker,
		startTime,
		endTime,
	)
}

func (s tradingContractServiceImpl) AcceptTradingRequest(ctx context.Context, id *big.Int) (*types.Transaction, error) {
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

	return s.tradingContract.AcceptTradingRequest(
		transactOpts,
		id,
	)
}

func (s tradingContractServiceImpl) DeclineTradingRequest(ctx context.Context, id *big.Int) (*types.Transaction, error) {
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

	return s.tradingContract.DeclineTradingRequest(
		transactOpts,
		id,
	)
}

func (s tradingContractServiceImpl) CreateTrade(ctx context.Context, negotiation *big.Int, broker common.Address) (*types.Transaction, error) {
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

	tx, err := s.tradingContract.CreateTrade(
		transactOpts,
		negotiation,
		broker,
	)
	if err != nil {
		_ = s.keyStore.Lock(account.Address)
		return nil, err
	}

	if err := s.keyStore.Lock(account.Address); err != nil {
		return nil, fmt.Errorf("lock account %s: %w", account.Address.Hex(), err)
	}
	return tx, err
}

func (s tradingContractServiceImpl) FindTradingRequestByIndex(ctx context.Context, index *big.Int) (*contracts.TradingRequest, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.tradingContract.FindTradingRequestByIndex(callOpts, index)
}

func (s tradingContractServiceImpl) FindTradeByIndex(ctx context.Context, index *big.Int) (*contracts.Trade, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.tradingContract.FindTradeByIndex(callOpts, index)
}

func (s tradingContractServiceImpl) FindTradingRequestById(ctx context.Context, id *big.Int) (*contracts.TradingRequest, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.tradingContract.FindTradingRequestById(callOpts, id)
}

func (s tradingContractServiceImpl) FindTradeById(ctx context.Context, id *big.Int) (*contracts.Trade, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.tradingContract.FindTradeById(callOpts, id)
}

func (s tradingContractServiceImpl) CountTradingRequests(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.tradingContract.CountTradingRequests(callOpts)
}

func (s tradingContractServiceImpl) CountTrades(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.tradingContract.CountTrades(callOpts)
}
