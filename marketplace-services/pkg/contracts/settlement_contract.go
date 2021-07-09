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

type Counter struct {
	Value *big.Int
	Set   bool
}

func CounterStructToCounter(result struct {
	Value *big.Int
	Set   bool
}) *Counter {
	return &Counter{
		Value: result.Value,
		Set:   result.Set,
	}
}

type Settlement struct {
	ActualCost *big.Int
	Provider   *big.Int
	Consumer   *big.Int
	Broker     *big.Int
}

func SettlementStructToSettlement(result struct {
	ActualCost *big.Int
	Provider   *big.Int
	Consumer   *big.Int
	Broker     *big.Int
}) *Settlement {
	return &Settlement{
		ActualCost: result.ActualCost,
		Provider:   result.Provider,
		Consumer:   result.Consumer,
		Broker:     result.Broker,
	}
}

type SettlementContract interface {
	Deposit(opts *bind.TransactOpts) (*types.Transaction, error)
	SettleTrade(opts *bind.TransactOpts, counter *big.Int) (*types.Transaction, error)
	ResolveDispute(opts *bind.TransactOpts, counter *big.Int) (*types.Transaction, error)
	ResolveTimeout(opts *bind.TransactOpts) (*types.Transaction, error)
	GetProviderCounter(opts *bind.CallOpts) (*Counter, error)
	GetConsumerCounter(opts *bind.CallOpts) (*Counter, error)
	GetBrokerCounter(opts *bind.CallOpts) (*Counter, error)
	GetSettlement(opts *bind.CallOpts) (*Settlement, error)
	WatchDepositedEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractDeposited, payee []common.Address) (event.Subscription, error)
	WatchSettledEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractSettled) (event.Subscription, error)
	WatchCounterSetEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractCounterSet, setter []common.Address) (event.Subscription, error)
	WatchDisputeEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractDispute) (event.Subscription, error)
}

type settlementContractImpl struct {
	binding *bindings.SettlementContract
}

func NewSettlementContractImpl(address common.Address, backend bind.ContractBackend) (*settlementContractImpl, error) {
	binding, err := bindings.NewSettlementContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new settlement contract binding %s: %w", address.Hex(), err)
	}
	return &settlementContractImpl{binding: binding}, err
}

func (s settlementContractImpl) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return s.binding.Deposit(opts)
}

func (s *settlementContractImpl) SettleTrade(opts *bind.TransactOpts, counter *big.Int) (*types.Transaction, error) {
	return s.binding.SettleTrade(opts, counter)
}

func (s *settlementContractImpl) ResolveDispute(opts *bind.TransactOpts, counter *big.Int) (*types.Transaction, error) {
	return s.binding.ResolveDispute(opts, counter)
}

func (s *settlementContractImpl) ResolveTimeout(opts *bind.TransactOpts) (*types.Transaction, error) {
	return s.binding.ResolveTimeout(opts)
}

func (s *settlementContractImpl) GetProviderCounter(opts *bind.CallOpts) (*Counter, error) {
	result, err := s.binding.GetProviderCounter(opts)
	return CounterStructToCounter(result), err
}

func (s *settlementContractImpl) GetConsumerCounter(opts *bind.CallOpts) (*Counter, error) {
	result, err := s.binding.GetConsumerCounter(opts)
	return CounterStructToCounter(result), err
}

func (s *settlementContractImpl) GetBrokerCounter(opts *bind.CallOpts) (*Counter, error) {
	result, err := s.binding.GetBrokerCounter(opts)
	return CounterStructToCounter(result), err
}

func (s *settlementContractImpl) GetSettlement(opts *bind.CallOpts) (*Settlement, error) {
	result, err := s.binding.GetSettlement(opts)
	return SettlementStructToSettlement(result), err
}

func (s settlementContractImpl) WatchDepositedEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractDeposited, payee []common.Address) (event.Subscription, error) {
	return s.binding.WatchDeposited(opts, sink, payee)
}

func (s *settlementContractImpl) WatchSettledEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractSettled) (event.Subscription, error) {
	return s.binding.WatchSettled(opts, sink)
}

func (s *settlementContractImpl) WatchCounterSetEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractCounterSet, setter []common.Address) (event.Subscription, error) {
	return s.binding.WatchCounterSet(opts, sink, setter)
}

func (s *settlementContractImpl) WatchDisputeEvent(opts *bind.WatchOpts, sink chan<- *bindings.SettlementContractDispute) (event.Subscription, error) {
	return s.binding.WatchDispute(opts, sink)
}
