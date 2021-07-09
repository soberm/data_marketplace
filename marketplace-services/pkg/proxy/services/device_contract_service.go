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

type DeviceContractService interface {
	CreateDevice(ctx context.Context, device *contracts.Device) (*types.Transaction, error)
	UpdateDevice(ctx context.Context, device *contracts.Device) (*types.Transaction, error)
	RemoveDevice(ctx context.Context, address common.Address) (*types.Transaction, error)
	FindDeviceByIndex(ctx context.Context, index *big.Int) (*contracts.Device, error)
	FindDeviceByAddress(ctx context.Context, address common.Address) (*contracts.Device, error)
	FindProductsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error)
	FindTradingRequestsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error)
	FindNegotiationRequestsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error)
	FindTradesOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error)
	FindNegotiationsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error)
	IsDeviceOwnedByUser(ctx context.Context, address common.Address, user common.Address) (bool, error)
	CountDevices(ctx context.Context) (*big.Int, error)
	ExistsDeviceByAddress(ctx context.Context, address common.Address) (bool, error)
	ExistsDeviceByAddressAndDeleted(ctx context.Context, address common.Address, deleted bool) (bool, error)
}

type deviceContractServiceImpl struct {
	logger         logrus.FieldLogger
	keyStore       *keystore.KeyStore
	walletService  WalletService
	deviceContract contracts.DeviceContract
}

func NewDeviceContractServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	deviceContract contracts.DeviceContract,
) *deviceContractServiceImpl {
	return &deviceContractServiceImpl{
		logger:         logger,
		walletService:  walletService,
		keyStore:       keyStore,
		deviceContract: deviceContract,
	}
}

func (s deviceContractServiceImpl) CreateDevice(ctx context.Context, device *contracts.Device) (*types.Transaction, error) {
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

	return s.deviceContract.CreateDevice(transactOpts, device)
}

func (s deviceContractServiceImpl) UpdateDevice(ctx context.Context, device *contracts.Device) (*types.Transaction, error) {
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

	return s.deviceContract.UpdateDevice(transactOpts, device)
}

func (s deviceContractServiceImpl) RemoveDevice(ctx context.Context, address common.Address) (*types.Transaction, error) {
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

	return s.deviceContract.RemoveDevice(transactOpts, address)
}

func (s deviceContractServiceImpl) FindDeviceByIndex(ctx context.Context, index *big.Int) (*contracts.Device, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.FindDeviceByIndex(callOpts, index)
}

func (s deviceContractServiceImpl) FindDeviceByAddress(ctx context.Context, address common.Address) (*contracts.Device, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.FindDeviceByAddress(callOpts, address)
}

func (s deviceContractServiceImpl) FindProductsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.FindProductsOfDeviceByAddress(callOpts, address)
}

func (s deviceContractServiceImpl) FindTradingRequestsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.FindTradingRequestsOfDeviceByAddress(callOpts, address)
}

func (s deviceContractServiceImpl) FindNegotiationRequestsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.FindNegotiationsOfDeviceByAddress(callOpts, address)
}

func (s deviceContractServiceImpl) FindTradesOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.FindTradesOfDeviceByAddress(callOpts, address)
}

func (s deviceContractServiceImpl) FindNegotiationsOfDeviceByAddress(ctx context.Context, address common.Address) ([]*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.FindNegotiationsOfDeviceByAddress(callOpts, address)
}

func (s deviceContractServiceImpl) IsDeviceOwnedByUser(ctx context.Context, address common.Address, user common.Address) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.IsDeviceOwnedByUser(callOpts, address, user)
}

func (s deviceContractServiceImpl) CountDevices(ctx context.Context) (*big.Int, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.CountDevices(callOpts)
}

func (s deviceContractServiceImpl) ExistsDeviceByAddress(ctx context.Context, address common.Address) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.ExistsDeviceByAddress(callOpts, address)
}

func (s deviceContractServiceImpl) ExistsDeviceByAddressAndDeleted(
	ctx context.Context,
	address common.Address,
	deleted bool,
) (bool, error) {
	w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
	if err != nil {
		return false, fmt.Errorf("find wallet of authenticated proxy account: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}
	return s.deviceContract.ExistsDeviceByAddressAndDeleted(callOpts, address, deleted)
}
