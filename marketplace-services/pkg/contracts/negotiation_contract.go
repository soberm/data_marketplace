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

type NegotiationRequest struct {
	Id       *big.Int
	Product  *big.Int
	Consumer common.Address
}

type Negotiation struct {
	Id              *big.Int
	Consumer        common.Address
	Product         *big.Int
	BiddingContract common.Address
}

func NegotiationRequestStructToNegotiationRequest(result struct {
	Id       *big.Int
	Product  *big.Int
	Consumer common.Address
}) *NegotiationRequest {
	return &NegotiationRequest{
		Id:       result.Id,
		Product:  result.Product,
		Consumer: result.Consumer,
	}
}

func NegotiationStructToNegotiation(result struct {
	Id              *big.Int
	Consumer        common.Address
	Product         *big.Int
	BiddingContract common.Address
}) *Negotiation {
	return &Negotiation{
		Id:              result.Id,
		Consumer:        result.Consumer,
		Product:         result.Product,
		BiddingContract: result.BiddingContract,
	}
}

type NegotiationContract interface {
	RequestNegotiation(opts *bind.TransactOpts, product *big.Int) (*types.Transaction, error)
	AcceptNegotiationRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	DeclineNegotiationRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	FindNegotiationRequestByIndex(opts *bind.CallOpts, index *big.Int) (*NegotiationRequest, error)
	FindNegotiationByIndex(opts *bind.CallOpts, index *big.Int) (*Negotiation, error)
	FindNegotiationRequestById(opts *bind.CallOpts, id *big.Int) (*NegotiationRequest, error)
	FindNegotiationById(opts *bind.CallOpts, id *big.Int) (*Negotiation, error)
	CountNegotiationRequests(opts *bind.CallOpts) (*big.Int, error)
	CountNegotiations(opts *bind.CallOpts) (*big.Int, error)
	ExistsNegotiationRequestById(opts *bind.CallOpts, id *big.Int) (bool, error)
	ExistsNegotiationById(opts *bind.CallOpts, id *big.Int) (bool, error)
	WatchRequestedNegotiationEvent(opts *bind.WatchOpts, sink chan<- *bindings.NegotiationContractRequestedNegotiation, requester []common.Address, product []*big.Int) (event.Subscription, error)
	WatchAcceptedNegotiationRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.NegotiationContractAcceptedNegotiationRequest, id []*big.Int) (event.Subscription, error)
	WatchDeclinedNegotiationRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.NegotiationContractDeclinedNegotiationRequest, id []*big.Int) (event.Subscription, error)
}

type negotiationContractImpl struct {
	binding *bindings.NegotiationContract
}

func NewNegotiationContractImpl(address common.Address, backend bind.ContractBackend) (*negotiationContractImpl, error) {
	binding, err := bindings.NewNegotiationContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new negotiation contract binding %s: %w", address.Hex(), err)
	}
	return &negotiationContractImpl{binding: binding}, err
}

func (n *negotiationContractImpl) RequestNegotiation(opts *bind.TransactOpts, product *big.Int) (*types.Transaction, error) {
	return n.binding.RequestNegotiation(opts, product)
}

func (n *negotiationContractImpl) AcceptNegotiationRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return n.binding.AcceptNegotiationRequest(opts, id)
}

func (n *negotiationContractImpl) DeclineNegotiationRequest(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return n.binding.DeclineNegotiationRequest(opts, id)
}

func (n *negotiationContractImpl) FindNegotiationRequestByIndex(opts *bind.CallOpts, index *big.Int) (*NegotiationRequest, error) {
	result, err := n.binding.FindNegotiationRequestByIndex(opts, index)
	return NegotiationRequestStructToNegotiationRequest(result), err
}

func (n *negotiationContractImpl) FindNegotiationByIndex(opts *bind.CallOpts, index *big.Int) (*Negotiation, error) {
	result, err := n.binding.FindNegotiationByIndex(opts, index)
	return NegotiationStructToNegotiation(result), err
}

func (n *negotiationContractImpl) FindNegotiationRequestById(opts *bind.CallOpts, id *big.Int) (*NegotiationRequest, error) {
	result, err := n.binding.FindNegotiationRequestById(opts, id)
	return NegotiationRequestStructToNegotiationRequest(result), err
}

func (n *negotiationContractImpl) FindNegotiationById(opts *bind.CallOpts, id *big.Int) (*Negotiation, error) {
	result, err := n.binding.FindNegotiationById(opts, id)
	return NegotiationStructToNegotiation(result), err
}

func (n *negotiationContractImpl) CountNegotiationRequests(opts *bind.CallOpts) (*big.Int, error) {
	return n.binding.CountNegotiationRequests(opts)
}

func (n *negotiationContractImpl) CountNegotiations(opts *bind.CallOpts) (*big.Int, error) {
	return n.binding.CountNegotiations(opts)
}

func (n *negotiationContractImpl) ExistsNegotiationRequestById(opts *bind.CallOpts, id *big.Int) (bool, error) {
	return n.binding.ExistsNegotiationRequestById(opts, id)
}

func (n *negotiationContractImpl) ExistsNegotiationById(opts *bind.CallOpts, id *big.Int) (bool, error) {
	return n.binding.ExistsNegotiationById(opts, id)
}

func (n *negotiationContractImpl) WatchRequestedNegotiationEvent(opts *bind.WatchOpts, sink chan<- *bindings.NegotiationContractRequestedNegotiation, requester []common.Address, product []*big.Int) (event.Subscription, error) {
	return n.binding.WatchRequestedNegotiation(opts, sink, requester, product)
}

func (n *negotiationContractImpl) WatchAcceptedNegotiationRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.NegotiationContractAcceptedNegotiationRequest, id []*big.Int) (event.Subscription, error) {
	return n.binding.WatchAcceptedNegotiationRequest(opts, sink, id)
}

func (n *negotiationContractImpl) WatchDeclinedNegotiationRequestEvent(opts *bind.WatchOpts, sink chan<- *bindings.NegotiationContractDeclinedNegotiationRequest, id []*big.Int) (event.Subscription, error) {
	return n.binding.WatchDeclinedNegotiationRequest(opts, sink, id)
}
