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

type TradingRequest struct {
	Id        *big.Int
	Product   *big.Int
	Cost      *big.Int
	StartTime *big.Int
	EndTime   *big.Int
	Consumer  common.Address
	Broker    common.Address
}

type Trade struct {
	Id                 *big.Int
	Provider           common.Address
	Consumer           common.Address
	Broker             common.Address
	Product            *big.Int
	StartTime          *big.Int
	EndTime            *big.Int
	Cost               *big.Int
	SettlementContract common.Address
}

func TradingRequestStructToTradingRequest(result struct {
	Id        *big.Int
	Product   *big.Int
	Cost      *big.Int
	StartTime *big.Int
	EndTime   *big.Int
	Consumer  common.Address
	Broker    common.Address
}) *TradingRequest {
	return &TradingRequest{
		Id:        result.Id,
		Product:   result.Product,
		Cost:      result.Cost,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
		Consumer:  result.Consumer,
		Broker:    result.Broker,
	}
}

func TradeStructToTrade(result struct {
	Id                 *big.Int
	Provider           common.Address
	Consumer           common.Address
	Broker             common.Address
	Product            *big.Int
	StartTime          *big.Int
	EndTime            *big.Int
	Cost               *big.Int
	SettlementContract common.Address
}) *Trade {
	return &Trade{
		Id:                 result.Id,
		Provider:           result.Provider,
		Consumer:           result.Consumer,
		Broker:             result.Broker,
		Product:            result.Product,
		StartTime:          result.StartTime,
		EndTime:            result.EndTime,
		Cost:               result.Cost,
		SettlementContract: result.SettlementContract,
	}
}

type TradingContract interface {
	RequestTrading(opts *bind.TransactOpts, product *big.Int, broker common.Address, startTime *big.Int, endTime *big.Int) (*types.Transaction, error)
	AcceptTradingRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	DeclineTradingRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	CreateTrade(opts *bind.TransactOpts, negotiation *big.Int, address common.Address) (*types.Transaction, error)
	FindTradingRequestByIndex(opts *bind.CallOpts, index *big.Int) (*TradingRequest, error)
	FindTradeByIndex(opts *bind.CallOpts, index *big.Int) (*Trade, error)
	FindTradingRequestById(opts *bind.CallOpts, id *big.Int) (*TradingRequest, error)
	FindTradeById(opts *bind.CallOpts, id *big.Int) (*Trade, error)
	CountTradingRequests(opts *bind.CallOpts) (*big.Int, error)
	CountTrades(opts *bind.CallOpts) (*big.Int, error)
	WatchRequestedTradingEvent(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractRequestedTrading, requester []common.Address, product []*big.Int) (event.Subscription, error)
	WatchAcceptedTradingRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractAcceptedTradingRequest, id []*big.Int) (event.Subscription, error)
	WatchDeclinedTradingRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractDeclinedTradingRequest, id []*big.Int) (event.Subscription, error)
	WatchCreatedTrade(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractCreatedTrade, broker []common.Address) (event.Subscription, error)
}

type tradingContractImpl struct {
	binding *bindings.TradingContract
}

func NewTradingContractImpl(address common.Address, backend bind.ContractBackend) (*tradingContractImpl, error) {
	binding, err := bindings.NewTradingContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new trading contract binding %s: %w", address.Hex(), err)
	}
	return &tradingContractImpl{binding: binding}, err
}

func (t *tradingContractImpl) RequestTrading(opts *bind.TransactOpts, product *big.Int, broker common.Address, startTime *big.Int, endTime *big.Int) (*types.Transaction, error) {
	return t.binding.RequestTrading(opts, product, broker, startTime, endTime)
}

func (t *tradingContractImpl) AcceptTradingRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return t.binding.AcceptTradingRequest(opts, id)
}

func (t *tradingContractImpl) DeclineTradingRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return t.binding.DeclineTradingRequest(opts, id)
}

func (t *tradingContractImpl) CreateTrade(opts *bind.TransactOpts, negotiation *big.Int, address common.Address) (*types.Transaction, error) {
	return t.binding.Create(opts, negotiation, address)
}

func (t *tradingContractImpl) FindTradingRequestByIndex(opts *bind.CallOpts, index *big.Int) (*TradingRequest, error) {
	result, err := t.binding.FindTradingRequestByIndex(opts, index)
	return TradingRequestStructToTradingRequest(result), err
}

func (t *tradingContractImpl) FindTradeByIndex(opts *bind.CallOpts, index *big.Int) (*Trade, error) {
	result, err := t.binding.FindTradeByIndex(opts, index)
	return TradeStructToTrade(result), err
}

func (t *tradingContractImpl) FindTradingRequestById(opts *bind.CallOpts, id *big.Int) (*TradingRequest, error) {
	result, err := t.binding.FindTradingRequestById(opts, id)
	return TradingRequestStructToTradingRequest(result), err
}

func (t *tradingContractImpl) FindTradeById(opts *bind.CallOpts, id *big.Int) (*Trade, error) {
	result, err := t.binding.FindTradeById(opts, id)
	return TradeStructToTrade(result), err
}

func (t *tradingContractImpl) CountTradingRequests(opts *bind.CallOpts) (*big.Int, error) {
	return t.binding.CountTradingRequests(opts)
}

func (t *tradingContractImpl) CountTrades(opts *bind.CallOpts) (*big.Int, error) {
	return t.binding.CountTrades(opts)
}

func (t *tradingContractImpl) WatchRequestedTradingEvent(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractRequestedTrading, requester []common.Address, product []*big.Int) (event.Subscription, error) {
	return t.binding.WatchRequestedTrading(opts, sink, requester, product)
}

func (t *tradingContractImpl) WatchAcceptedTradingRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractAcceptedTradingRequest, id []*big.Int) (event.Subscription, error) {
	return t.binding.WatchAcceptedTradingRequest(opts, sink, id)
}

func (t *tradingContractImpl) WatchDeclinedTradingRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractDeclinedTradingRequest, id []*big.Int) (event.Subscription, error) {
	return t.binding.WatchDeclinedTradingRequest(opts, sink, id)
}

func (t *tradingContractImpl) WatchCreatedTrade(opts *bind.WatchOpts, sink chan<- *bindings.TradingContractCreatedTrade, broker []common.Address) (event.Subscription, error) {
	return t.binding.WatchCreatedTrade(opts, sink, broker)
}
