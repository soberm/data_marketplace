package contracts

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"marketplace-services/pkg/contracts/bindings"
	"math/big"
)

type Device struct {
	Addr        common.Address
	User        string
	Name        string
	Description string
	PublicKey   []byte
	Rating      *big.Int
	Deleted     bool
}

func DeviceStructToDevice(result struct {
	Addr        common.Address
	User        common.Address
	Name        string
	Description string
	PublicKey   []byte
	Rating      *big.Int
	Deleted     bool
}) *Device {
	return &Device{
		Addr:        result.Addr,
		User:        result.User.Hex(),
		Name:        result.Name,
		Description: result.Description,
		PublicKey:   result.PublicKey,
		Rating:      result.Rating,
		Deleted:     result.Deleted,
	}
}

type DeviceContract interface {
	CreateDevice(opts *bind.TransactOpts, device *Device) (*types.Transaction, error)
	UpdateDevice(opts *bind.TransactOpts, device *Device) (*types.Transaction, error)
	RemoveDevice(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error)
	FindDeviceByIndex(opts *bind.CallOpts, index *big.Int) (*Device, error)
	FindDeviceByAddress(opts *bind.CallOpts, address common.Address) (*Device, error)
	FindProductsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error)
	FindTradingRequestsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error)
	FindNegotiationRequestsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error)
	FindTradesOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error)
	FindNegotiationsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error)
	IsDeviceOwnedByUser(opts *bind.CallOpts, address common.Address, user common.Address) (bool, error)
	CountDevices(opts *bind.CallOpts) (*big.Int, error)
	ExistsDeviceByAddress(opts *bind.CallOpts, address common.Address) (bool, error)
	ExistsDeviceByAddressAndDeleted(opts *bind.CallOpts, address common.Address, deleted bool) (bool, error)
	WatchCreatedDeviceEvent(opts *bind.WatchOpts, sink chan<- *bindings.DeviceContractCreatedDevice, addr []common.Address) (event.Subscription, error)
	WatchUpdatedDeviceEvent(opts *bind.WatchOpts, sink chan<- *bindings.DeviceContractUpdatedDevice, addr []common.Address) (event.Subscription, error)
	WatchRemovedDeviceEvent(opts *bind.WatchOpts, sink chan<- *bindings.DeviceContractRemovedDevice, addr []common.Address) (event.Subscription, error)
}

type deviceContractImpl struct {
	binding *bindings.DeviceContract
}

func NewDeviceContractImpl(address common.Address, backend bind.ContractBackend) (*deviceContractImpl, error) {
	binding, err := bindings.NewDeviceContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new device contract binding %s: %w", address.Hex(), err)
	}
	return &deviceContractImpl{binding: binding}, err
}

func (d *deviceContractImpl) CreateDevice(opts *bind.TransactOpts, device *Device) (*types.Transaction, error) {
	return d.binding.Create(opts, device.Addr, device.Name, device.Description, device.PublicKey)
}

func (d *deviceContractImpl) UpdateDevice(opts *bind.TransactOpts, device *Device) (*types.Transaction, error) {
	return d.binding.Update(opts, device.Addr, device.Name, device.Description, device.PublicKey)
}

func (d *deviceContractImpl) RemoveDevice(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return d.binding.Remove(opts, addr)
}

func (d *deviceContractImpl) FindDeviceByIndex(opts *bind.CallOpts, index *big.Int) (*Device, error) {
	result, err := d.binding.FindByIndex(opts, index)
	return DeviceStructToDevice(result), err
}

func (d *deviceContractImpl) FindDeviceByAddress(opts *bind.CallOpts, address common.Address) (*Device, error) {
	result, err := d.binding.FindByAddress(opts, address)
	return DeviceStructToDevice(result), err
}

func (d *deviceContractImpl) FindProductsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error) {
	return d.binding.FindProductsByAddress(opts, address)
}

func (d *deviceContractImpl) FindTradingRequestsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error) {
	return d.binding.FindTradingRequestsByAddress(opts, address)
}

func (d *deviceContractImpl) FindNegotiationRequestsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error) {
	return d.binding.FindNegotiationRequestsByAddress(opts, address)
}

func (d *deviceContractImpl) FindTradesOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error) {
	return d.binding.FindTradesByAddress(opts, address)
}

func (d *deviceContractImpl) FindNegotiationsOfDeviceByAddress(opts *bind.CallOpts, address common.Address) ([]*big.Int, error) {
	return d.binding.FindNegotiationsByAddress(opts, address)
}

func (d *deviceContractImpl) IsDeviceOwnedByUser(opts *bind.CallOpts, address common.Address, user common.Address) (bool, error) {
	return d.binding.IsOwnedByUser(opts, address, user)
}

func (d *deviceContractImpl) CountDevices(opts *bind.CallOpts) (*big.Int, error) {
	return d.binding.Count(opts)
}

func (d *deviceContractImpl) ExistsDeviceByAddress(opts *bind.CallOpts, address common.Address) (bool, error) {
	return d.binding.ExistsByAddress(opts, address)
}

func (d *deviceContractImpl) ExistsDeviceByAddressAndDeleted(opts *bind.CallOpts, address common.Address, deleted bool) (bool, error) {
	return d.binding.ExistsByAddressAndDeleted(opts, address, deleted)
}

func (d *deviceContractImpl) WatchCreatedDeviceEvent(opts *bind.WatchOpts, sink chan<- *bindings.DeviceContractCreatedDevice, addr []common.Address) (event.Subscription, error) {
	return d.binding.WatchCreatedDevice(opts, sink, addr)
}

func (d *deviceContractImpl) WatchUpdatedDeviceEvent(opts *bind.WatchOpts, sink chan<- *bindings.DeviceContractUpdatedDevice, addr []common.Address) (event.Subscription, error) {
	return d.binding.WatchUpdatedDevice(opts, sink, addr)
}

func (d *deviceContractImpl) WatchRemovedDeviceEvent(opts *bind.WatchOpts, sink chan<- *bindings.DeviceContractRemovedDevice, addr []common.Address) (event.Subscription, error) {
	return d.binding.WatchRemovedDevice(opts, sink, addr)
}
