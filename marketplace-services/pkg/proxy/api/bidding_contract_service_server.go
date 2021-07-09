package api

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/proxy/services"
	"math/big"
)

type biddingContractServiceServer struct {
	UnimplementedBiddingContractServiceServer
	logger        logrus.FieldLogger
	walletService services.WalletService
	keyStore      *keystore.KeyStore
	ethClient     *ethclient.Client
}

func NewBiddingContractServiceServer(
	logger logrus.FieldLogger,
	walletService services.WalletService,
	keyStore *keystore.KeyStore,
	ethClient *ethclient.Client,
) *biddingContractServiceServer {
	return &biddingContractServiceServer{
		logger:        logger,
		walletService: walletService,
		keyStore:      keyStore,
		ethClient:     ethClient,
	}
}

func (s *biddingContractServiceServer) MakeBid(ctx context.Context, req *MakeBidRequest) (*MakeBidResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &MakeBidResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	tx, err := service.MakeABid(ctx, BidFromGrpcBid(req.Bid))
	if err != nil {
		return &MakeBidResponse{}, err
	}
	return &MakeBidResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *biddingContractServiceServer) AcceptLastBid(
	ctx context.Context,
	req *AcceptLastBidRequest,
) (*AcceptLastBidResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &AcceptLastBidResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	tx, err := service.AcceptLastBid(ctx)
	if err != nil {
		return &AcceptLastBidResponse{}, err
	}
	return &AcceptLastBidResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *biddingContractServiceServer) CancelBidding(
	ctx context.Context,
	req *CancelBiddingRequest,
) (*CancelBiddingResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &CancelBiddingResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	tx, err := service.CancelBidding(ctx)
	if err != nil {
		return &CancelBiddingResponse{}, err
	}
	return &CancelBiddingResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *biddingContractServiceServer) FindBidByIndex(
	ctx context.Context,
	req *FindBidByIndexRequest,
) (*FindBidByIndexResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &FindBidByIndexResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	b, err := service.FindBidByIndex(ctx, big.NewInt(int64(req.Index)))
	if err != nil {
		return &FindBidByIndexResponse{}, err
	}
	return &FindBidByIndexResponse{Bid: BidToGrpcBid(b)}, err
}

func (s *biddingContractServiceServer) FindLastBid(
	ctx context.Context,
	req *FindLastBidRequest,
) (*FindLastBidResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &FindLastBidResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	b, err := service.FindLastBid(ctx)
	if err != nil {
		return &FindLastBidResponse{}, err
	}
	return &FindLastBidResponse{Bid: BidToGrpcBid(b)}, err
}
func (s *biddingContractServiceServer) CountBids(
	ctx context.Context,
	req *CountBidsRequest,
) (*CountBidsResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &CountBidsResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	c, err := service.CountBids(ctx)
	if err != nil {
		return &CountBidsResponse{}, err
	}
	return &CountBidsResponse{Counter: c.Uint64()}, err
}
func (s *biddingContractServiceServer) IsLastBidAccepted(
	ctx context.Context,
	req *IsLastBidAcceptedRequest,
) (*IsLastBidAcceptedResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &IsLastBidAcceptedResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	accepted, err := service.IsLastBidAccepted(ctx)
	if err != nil {
		return &IsLastBidAcceptedResponse{}, err
	}
	return &IsLastBidAcceptedResponse{Accepted: accepted}, err
}
func (s *biddingContractServiceServer) IsBiddingCanceled(
	ctx context.Context,
	req *IsBiddingCanceledRequest,
) (*IsBiddingCanceledResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &IsBiddingCanceledResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	accepted, err := service.IsLastBidAccepted(ctx)
	if err != nil {
		return &IsBiddingCanceledResponse{}, err
	}
	return &IsBiddingCanceledResponse{Canceled: accepted}, err
}
func (s *biddingContractServiceServer) IsBiddingActive(
	ctx context.Context,
	req *IsBiddingActiveRequest,
) (*IsBiddingActiveResponse, error) {
	contract, err := contracts.NewBiddingContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &IsBiddingActiveResponse{}, err
	}
	service := services.NewBiddingContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	active, err := service.IsLastBidAccepted(ctx)
	if err != nil {
		return &IsBiddingActiveResponse{}, err
	}
	return &IsBiddingActiveResponse{Active: active}, err
}
