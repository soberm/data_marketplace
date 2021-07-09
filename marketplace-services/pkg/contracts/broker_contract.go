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

type Location uint8

const (
	LocationBR Location = iota
	LocationEUNE
	LocationEUW
	LocationLAN
	LocationLAS
	LocationNA
	LocationOCE
	LocationRU
	LocationTR
	LocationJP
	LocationPH
	LocationSG
	LocationTW
	LocationVN
	LocationTH
	LocationKR
	LocationCN
)

type Broker struct {
	Addr     common.Address
	User     common.Address
	Name     string
	HostAddr string
	Location Location
	Trades   []*big.Int
	Deleted  bool
}

func BrokerStructToBroker(result struct {
	Addr     common.Address
	User     common.Address
	Name     string
	HostAddr string
	Location uint8
	Trades   []*big.Int
	Deleted  bool
}) *Broker {
	return &Broker{
		Addr:     result.Addr,
		User:     result.User,
		Name:     result.Name,
		HostAddr: result.HostAddr,
		Location: Location(result.Location),
		Trades:   result.Trades,
		Deleted:  result.Deleted,
	}
}

type BrokerContract interface {
	CreateBroker(opts *bind.TransactOpts, broker *Broker) (*types.Transaction, error)
	UpdateBroker(opts *bind.TransactOpts, broker *Broker) (*types.Transaction, error)
	RemoveBroker(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error)
	FindBrokerByIndex(opts *bind.CallOpts, index *big.Int) (*Broker, error)
	FindBrokerByAddress(opts *bind.CallOpts, addr common.Address) (*Broker, error)
	CountBrokers(opts *bind.CallOpts) (*big.Int, error)
	ExistsBrokerByAddress(opts *bind.CallOpts, addr common.Address) (bool, error)
	ExistsBrokerByAddressAndDeleted(opts *bind.CallOpts, addr common.Address, deleted bool) (bool, error)
	WatchCreatedBrokerEvent(opts *bind.WatchOpts, sink chan<- *bindings.BrokerContractCreatedBroker, addr []common.Address) (event.Subscription, error)
	WatchUpdatedBrokerEvent(opts *bind.WatchOpts, sink chan<- *bindings.BrokerContractUpdatedBroker, addr []common.Address) (event.Subscription, error)
	WatchRemovedBrokerEvent(opts *bind.WatchOpts, sink chan<- *bindings.BrokerContractRemovedBroker, addr []common.Address) (event.Subscription, error)
}

type brokerContractImpl struct {
	binding *bindings.BrokerContract
}

func NewBrokerContractImpl(address common.Address, backend bind.ContractBackend) (*brokerContractImpl, error) {
	binding, err := bindings.NewBrokerContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new broker contract binding %s: %w", address.Hex(), err)
	}
	return &brokerContractImpl{binding: binding}, err
}

func (b *brokerContractImpl) CreateBroker(opts *bind.TransactOpts, broker *Broker) (*types.Transaction, error) {
	return b.binding.Create(opts, broker.Addr, broker.Name, broker.HostAddr, uint8(broker.Location))
}

func (b *brokerContractImpl) UpdateBroker(opts *bind.TransactOpts, broker *Broker) (*types.Transaction, error) {
	return b.binding.Update(opts, broker.Addr, broker.Name, broker.HostAddr, uint8(broker.Location))
}

func (b *brokerContractImpl) RemoveBroker(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return b.binding.Remove(opts, addr)
}

func (b *brokerContractImpl) FindBrokerByIndex(opts *bind.CallOpts, index *big.Int) (*Broker, error) {
	result, err := b.binding.FindByIndex(opts, index)
	return BrokerStructToBroker(result), err
}

func (b *brokerContractImpl) FindBrokerByAddress(opts *bind.CallOpts, addr common.Address) (*Broker, error) {
	result, err := b.binding.FindByAddress(opts, addr)
	return BrokerStructToBroker(result), err
}

func (b *brokerContractImpl) CountBrokers(opts *bind.CallOpts) (*big.Int, error) {
	return b.binding.Count(opts)
}

func (b *brokerContractImpl) ExistsBrokerByAddress(opts *bind.CallOpts, addr common.Address) (bool, error) {
	return b.ExistsBrokerByAddress(opts, addr)
}

func (b *brokerContractImpl) ExistsBrokerByAddressAndDeleted(opts *bind.CallOpts, addr common.Address, deleted bool) (bool, error) {
	return b.ExistsBrokerByAddressAndDeleted(opts, addr, deleted)
}

func (b *brokerContractImpl) WatchCreatedBrokerEvent(opts *bind.WatchOpts, sink chan<- *bindings.BrokerContractCreatedBroker, addr []common.Address) (event.Subscription, error) {
	return b.binding.WatchCreatedBroker(opts, sink, addr)
}

func (b *brokerContractImpl) WatchUpdatedBrokerEvent(opts *bind.WatchOpts, sink chan<- *bindings.BrokerContractUpdatedBroker, addr []common.Address) (event.Subscription, error) {
	return b.binding.WatchUpdatedBroker(opts, sink, addr)
}

func (b *brokerContractImpl) WatchRemovedBrokerEvent(opts *bind.WatchOpts, sink chan<- *bindings.BrokerContractRemovedBroker, addr []common.Address) (event.Subscription, error) {
	return b.binding.WatchRemovedBroker(opts, sink, addr)
}
