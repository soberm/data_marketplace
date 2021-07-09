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

type NegotiationContractService interface {
	RequestNegotiation(ctx context.Context, product *big.Int) (*types.Transaction, error)
	AcceptNegotiationRequest(ctx context.Context, id *big.Int) (*types.Transaction, error)
	DeclineNegotiationRequest(ctx context.Context, id *big.Int) (*types.Transaction, error)
	FindNegotiationRequestByIndex(ctx context.Context, index *big.Int) (*contracts.NegotiationRequest, error)
	FindNegotiationByIndex(ctx context.Context, index *big.Int) (*contracts.Negotiation, error)
	FindNegotiationRequestById(ctx context.Context, id *big.Int) (*contracts.NegotiationRequest, error)
	FindNegotiationById(ctx context.Context, id *big.Int) (*contracts.Negotiation, error)
	CountNegotiationRequests(ctx context.Context) (*big.Int, error)
	CountNegotiations(ctx context.Context) (*big.Int, error)
	ExistsNegotiationRequestById(ctx context.Context, id *big.Int) (bool, error)
	ExistsNegotiationById(ctx context.Context, id *big.Int) (bool, error)
}

type negotiationContractServiceImpl struct {
	logger              logrus.FieldLogger
	keyStore            *keystore.KeyStore
	walletService       WalletService
	negotiationContract contracts.NegotiationContract
}

func NewNegotiationContractServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	negotiationContract contracts.NegotiationContract,
) *negotiationContractServiceImpl {
	return &negotiationContractServiceImpl{
		logger:              logger,
		walletService:       walletService,
		keyStore:            keyStore,
		negotiationContract: negotiationContract,
	}
}

func (s negotiationContractServiceImpl) RequestNegotiation(ctx context.Context, product *big.Int) (*types.Transaction, error) {
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

	return s.negotiationContract.RequestNegotiation(transactOpts, product)
}

func (s negotiationContractServiceImpl) AcceptNegotiationRequest(ctx context.Context, id *big.Int) (*types.Transaction, error) {
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

	return s.negotiationContract.AcceptNegotiationRequest(transactOpts, id)
}

func (s negotiationContractServiceImpl) DeclineNegotiationRequest(ctx context.Context, id *big.Int) (*types.Transaction, error) {
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

	return s.negotiationContract.DeclineNegotiationRequest(transactOpts, id)
}

func (s negotiationContractServiceImpl) FindNegotiationRequestByIndex(ctx context.Context, index *big.Int) (*contracts.NegotiationRequest, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.FindNegotiationRequestByIndex(callOpts, index)
}

func (s negotiationContractServiceImpl) FindNegotiationByIndex(ctx context.Context, index *big.Int) (*contracts.Negotiation, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.FindNegotiationByIndex(callOpts, index)
}

func (s negotiationContractServiceImpl) FindNegotiationRequestById(ctx context.Context, id *big.Int) (*contracts.NegotiationRequest, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.FindNegotiationRequestById(callOpts, id)
}

func (s negotiationContractServiceImpl) FindNegotiationById(ctx context.Context, id *big.Int) (*contracts.Negotiation, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.FindNegotiationById(callOpts, id)
}

func (s negotiationContractServiceImpl) CountNegotiationRequests(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.CountNegotiationRequests(callOpts)
}

func (s negotiationContractServiceImpl) CountNegotiations(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.CountNegotiations(callOpts)
}

func (s negotiationContractServiceImpl) ExistsNegotiationRequestById(ctx context.Context, id *big.Int) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.ExistsNegotiationRequestById(callOpts, id)
}

func (s negotiationContractServiceImpl) ExistsNegotiationById(ctx context.Context, id *big.Int) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.negotiationContract.ExistsNegotiationById(callOpts, id)
}
