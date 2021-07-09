package contracts

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"marketplace-services/pkg/contracts/bindings"
	"math/big"
)

type Bid struct {
	Price     *big.Int
	StartTime *big.Int
	EndTime   *big.Int
}

func BidStructToBid(result struct {
	Price     *big.Int
	StartTime *big.Int
	EndTime   *big.Int
}) *Bid {
	return &Bid{
		Price:     result.Price,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
	}
}

type BiddingContract interface {
	MakeABid(opts *bind.TransactOpts, bid *Bid) (*types.Transaction, error)
	AcceptLastBid(opts *bind.TransactOpts) (*types.Transaction, error)
	CancelBidding(opts *bind.TransactOpts) (*types.Transaction, error)
	FindBidByIndex(opts *bind.CallOpts, index *big.Int) (*Bid, error)
	FindLastBid(opts *bind.CallOpts) (*Bid, error)
	CountBids(opts *bind.CallOpts) (*big.Int, error)
	IsLastBidAccepted(opts *bind.CallOpts) (bool, error)
	IsBiddingCanceled(opts *bind.CallOpts) (bool, error)
	IsBiddingActive(opts *bind.CallOpts) (bool, error)
}

type biddingContractImpl struct {
	binding *bindings.BiddingContract
}

func NewBiddingContractImpl(address common.Address, backend bind.ContractBackend) (*biddingContractImpl, error) {
	binding, err := bindings.NewBiddingContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new bidding contract binding %s: %w", address.Hex(), err)
	}
	return &biddingContractImpl{binding: binding}, err
}

func (b biddingContractImpl) MakeABid(opts *bind.TransactOpts, bid *Bid) (*types.Transaction, error) {
	return b.binding.MakeBid(opts, bid.Price, bid.StartTime, bid.EndTime)
}

func (b biddingContractImpl) AcceptLastBid(opts *bind.TransactOpts) (*types.Transaction, error) {
	return b.binding.Accept(opts)
}

func (b biddingContractImpl) CancelBidding(opts *bind.TransactOpts) (*types.Transaction, error) {
	return b.binding.Cancel(opts)
}

func (b biddingContractImpl) FindBidByIndex(opts *bind.CallOpts, index *big.Int) (*Bid, error) {
	result, err := b.binding.FindByIndex(opts, index)
	return BidStructToBid(result), err
}

func (b biddingContractImpl) FindLastBid(opts *bind.CallOpts) (*Bid, error) {
	result, err := b.binding.LastBid(opts)
	return BidStructToBid(result), err
}

func (b biddingContractImpl) CountBids(opts *bind.CallOpts) (*big.Int, error) {
	return b.binding.Count(opts)
}

func (b biddingContractImpl) IsLastBidAccepted(opts *bind.CallOpts) (bool, error) {
	return b.binding.IsAccepted(opts)
}

func (b biddingContractImpl) IsBiddingCanceled(opts *bind.CallOpts) (bool, error) {
	return b.binding.IsCanceled(opts)
}

func (b biddingContractImpl) IsBiddingActive(opts *bind.CallOpts) (bool, error) {
	return b.binding.IsActive(opts)
}
