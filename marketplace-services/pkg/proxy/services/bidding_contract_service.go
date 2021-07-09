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

type BiddingContractService interface {
	MakeABid(ctx context.Context, bid *contracts.Bid) (*types.Transaction, error)
	AcceptLastBid(ctx context.Context) (*types.Transaction, error)
	CancelBidding(ctx context.Context) (*types.Transaction, error)
	FindBidByIndex(ctx context.Context, index *big.Int) (*contracts.Bid, error)
	FindLastBid(ctx context.Context) (*contracts.Bid, error)
	CountBids(ctx context.Context) (*big.Int, error)
	IsLastBidAccepted(ctx context.Context) (bool, error)
	IsBiddingCanceled(ctx context.Context) (bool, error)
	IsBiddingActive(ctx context.Context) (bool, error)
}

type biddingContractServiceImpl struct {
	logger          logrus.FieldLogger
	keyStore        *keystore.KeyStore
	walletService   WalletService
	biddingContract contracts.BiddingContract
}

func NewBiddingContractServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	biddingContract contracts.BiddingContract,
) *biddingContractServiceImpl {
	return &biddingContractServiceImpl{
		logger:          logger,
		walletService:   walletService,
		keyStore:        keyStore,
		biddingContract: biddingContract,
	}
}

func (s biddingContractServiceImpl) MakeABid(ctx context.Context, bid *contracts.Bid) (*types.Transaction, error) {
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

	return s.biddingContract.MakeABid(transactOpts, bid)
}

func (s biddingContractServiceImpl) AcceptLastBid(ctx context.Context) (*types.Transaction, error) {
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

	return s.biddingContract.AcceptLastBid(transactOpts)
}

func (s biddingContractServiceImpl) CancelBidding(ctx context.Context) (*types.Transaction, error) {
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

	return s.biddingContract.CancelBidding(transactOpts)
}

func (s biddingContractServiceImpl) FindBidByIndex(ctx context.Context, index *big.Int) (*contracts.Bid, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.biddingContract.FindBidByIndex(callOpts, index)
}

func (s biddingContractServiceImpl) FindLastBid(ctx context.Context) (*contracts.Bid, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.biddingContract.FindLastBid(callOpts)
}

func (s biddingContractServiceImpl) CountBids(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.biddingContract.CountBids(callOpts)
}

func (s biddingContractServiceImpl) IsLastBidAccepted(ctx context.Context) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.biddingContract.IsLastBidAccepted(callOpts)
}

func (s biddingContractServiceImpl) IsBiddingCanceled(ctx context.Context) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.biddingContract.IsBiddingCanceled(callOpts)
}

func (s biddingContractServiceImpl) IsBiddingActive(ctx context.Context) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.biddingContract.IsBiddingActive(callOpts)
}
