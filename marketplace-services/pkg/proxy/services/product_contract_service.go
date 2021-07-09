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

type ProductContractService interface {
	CreateProduct(ctx context.Context, product *contracts.Product) (*types.Transaction, error)
	UpdateProduct(ctx context.Context, product *contracts.Product) (*types.Transaction, error)
	RemoveProduct(ctx context.Context, id *big.Int) (*types.Transaction, error)
	FindProductByIndex(ctx context.Context, index *big.Int) (*contracts.Product, error)
	FindProductById(ctx context.Context, id *big.Int) (*contracts.Product, error)
	FindTradingRequestsOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error)
	FindNegotiationRequestsOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error)
	FindTradesOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error)
	FindNegotiationsOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error)
	FindCostOfProductById(ctx context.Context, id *big.Int) (*big.Int, error)
	IsProductOwnedByDevice(ctx context.Context, id *big.Int, device common.Address) (bool, error)
	CountProducts(ctx context.Context) (*big.Int, error)
	ExistsProductById(ctx context.Context, id *big.Int) (bool, error)
	ExistsProductByIdAndDeleted(ctx context.Context, id *big.Int, deleted bool) (bool, error)
}

type productContractServiceImpl struct {
	logger          logrus.FieldLogger
	keyStore        *keystore.KeyStore
	walletService   WalletService
	productContract contracts.ProductContract
}

func NewProductContractService(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	productContract contracts.ProductContract,
) *productContractServiceImpl {
	return &productContractServiceImpl{logger: logger, walletService: walletService, keyStore: keyStore, productContract: productContract}
}

func (s productContractServiceImpl) CreateProduct(ctx context.Context, product *contracts.Product) (*types.Transaction, error) {
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

	return s.productContract.CreateProduct(transactOpts, product)
}

func (s productContractServiceImpl) UpdateProduct(ctx context.Context, product *contracts.Product) (*types.Transaction, error) {
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

	return s.productContract.UpdateProduct(transactOpts, product)
}

func (s productContractServiceImpl) RemoveProduct(ctx context.Context, id *big.Int) (*types.Transaction, error) {
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

	return s.productContract.RemoveProduct(transactOpts, id)
}

func (s productContractServiceImpl) FindProductByIndex(ctx context.Context, index *big.Int) (*contracts.Product, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.FindProductByIndex(callOpts, index)
}

func (s productContractServiceImpl) FindProductById(ctx context.Context, id *big.Int) (*contracts.Product, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.FindProductById(callOpts, id)
}

func (s productContractServiceImpl) FindTradingRequestsOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.FindTradingRequestsOfProductById(callOpts, id)
}

func (s productContractServiceImpl) FindNegotiationRequestsOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.FindNegotiationRequestsOfProductById(callOpts, id)
}

func (s productContractServiceImpl) FindTradesOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.FindTradesOfProductById(callOpts, id)
}

func (s productContractServiceImpl) FindNegotiationsOfProductById(ctx context.Context, id *big.Int) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.FindNegotiationsOfProductById(callOpts, id)
}

func (s productContractServiceImpl) FindCostOfProductById(ctx context.Context, id *big.Int) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.FindCostOfProductById(callOpts, id)
}

func (s productContractServiceImpl) IsProductOwnedByDevice(ctx context.Context, id *big.Int, device common.Address) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.IsProductOwnedByDevice(callOpts, id, device)
}

func (s productContractServiceImpl) CountProducts(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.CountProducts(callOpts)
}

func (s productContractServiceImpl) ExistsProductById(ctx context.Context, id *big.Int) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.ExistsProductById(callOpts, id)
}

func (s productContractServiceImpl) ExistsProductByIdAndDeleted(ctx context.Context, id *big.Int, deleted bool) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.productContract.ExistsProductByIdAndDeleted(callOpts, id, deleted)
}
