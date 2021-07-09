package api

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/contracts/bindings"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/services"
	"math/big"
)

type tradingContractServiceServer struct {
	UnimplementedTradingContractServiceServer
	tradingContractService services.TradingContractService
	tradingContract        contracts.TradingContract
}

func NewTradingContractServiceServer(
	tradingContractService services.TradingContractService,
	tradingContract contracts.TradingContract,
) *tradingContractServiceServer {
	return &tradingContractServiceServer{
		tradingContractService: tradingContractService,
		tradingContract:        tradingContract,
	}
}

func (s *tradingContractServiceServer) RequestTrading(
	ctx context.Context,
	req *RequestTradingRequest,
) (*RequestTradingResponse, error) {
	tx, err := s.tradingContractService.RequestTrading(
		ctx,
		big.NewInt(int64(req.Product)),
		common.HexToAddress(req.Broker),
		big.NewInt(int64(req.StartTime)),
		big.NewInt(int64(req.EndTime)),
	)
	if err != nil {
		return &RequestTradingResponse{}, err
	}
	return &RequestTradingResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *tradingContractServiceServer) AcceptTradingRequest(
	ctx context.Context,
	req *AcceptTradingRequestRequest,
) (*AcceptTradingRequestResponse, error) {
	tx, err := s.tradingContractService.AcceptTradingRequest(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &AcceptTradingRequestResponse{}, err
	}
	return &AcceptTradingRequestResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *tradingContractServiceServer) DeclineTradingRequest(
	ctx context.Context,
	req *DeclineTradingRequestRequest,
) (*DeclineTradingRequestResponse, error) {
	tx, err := s.tradingContractService.AcceptTradingRequest(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &DeclineTradingRequestResponse{}, err
	}
	return &DeclineTradingRequestResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *tradingContractServiceServer) CreateTrade(
	ctx context.Context,
	req *CreateTradeRequest,
) (*CreateTradeResponse, error) {
	tx, err := s.tradingContractService.CreateTrade(ctx, big.NewInt(int64(req.Negotiation)), common.HexToAddress(req.Broker))
	if err != nil {
		return &CreateTradeResponse{}, err
	}
	return &CreateTradeResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *tradingContractServiceServer) FindTradingRequestByIndex(
	ctx context.Context,
	req *FindTradingRequestByIndexRequest,
) (*FindTradingRequestByIndexResponse, error) {
	tr, err := s.tradingContractService.FindTradingRequestByIndex(ctx, big.NewInt(int64(req.Index)))
	if err != nil {
		return &FindTradingRequestByIndexResponse{}, err
	}
	return &FindTradingRequestByIndexResponse{Request: TradingRequestToGrpcTradingRequest(tr)}, err
}
func (s *tradingContractServiceServer) FindTradeByIndex(
	ctx context.Context,
	req *FindTradeByIndexRequest,
) (*FindTradeByIndexResponse, error) {
	t, err := s.tradingContractService.FindTradeByIndex(ctx, big.NewInt(int64(req.Index)))
	if err != nil {
		return &FindTradeByIndexResponse{}, err
	}
	return &FindTradeByIndexResponse{Trade: TradeToGrpcTrade(t)}, err
}
func (s *tradingContractServiceServer) FindTradingRequestById(
	ctx context.Context,
	req *FindTradingRequestByIdRequest,
) (*FindTradingRequestByIdResponse, error) {
	tr, err := s.tradingContractService.FindTradingRequestById(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &FindTradingRequestByIdResponse{}, err
	}
	return &FindTradingRequestByIdResponse{Request: TradingRequestToGrpcTradingRequest(tr)}, err
}
func (s *tradingContractServiceServer) FindTradeById(
	ctx context.Context,
	req *FindTradeByIdRequest,
) (*FindTradeByIdResponse, error) {
	t, err := s.tradingContractService.FindTradeById(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &FindTradeByIdResponse{}, err
	}
	return &FindTradeByIdResponse{Trade: TradeToGrpcTrade(t)}, err
}
func (s *tradingContractServiceServer) CountTradingRequests(
	ctx context.Context,
	_ *CountTradingRequestsRequest,
) (*CountTradingRequestsResponse, error) {
	c, err := s.tradingContractService.CountTradingRequests(ctx)
	if err != nil {
		return &CountTradingRequestsResponse{}, err
	}
	return &CountTradingRequestsResponse{Counter: c.Uint64()}, err
}
func (s *tradingContractServiceServer) CountTrades(
	ctx context.Context,
	_ *CountTradesRequest,
) (*CountTradesResponse, error) {
	c, err := s.tradingContractService.CountTrades(ctx)
	if err != nil {
		return &CountTradesResponse{}, err
	}
	return &CountTradesResponse{Counter: c.Uint64()}, err
}
func (s *tradingContractServiceServer) WatchRequestedTradingEvent(
	req *WatchRequestedTradingEventRequest,
	stream TradingContractService_WatchRequestedTradingEventServer,
) error {
	sink := make(chan *bindings.TradingContractRequestedTrading)
	defer close(sink)

	sub, err := s.tradingContract.WatchRequestedTradingEvent(
		&bind.WatchOpts{
			Context: stream.Context(),
		},
		sink,
		contracts.HexToAddresses(req.Requesters),
		contracts.UInt64ToBigInt(req.Products),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			if err := stream.Send(&WatchRequestedTradingEventResponse{
				Event: &domain.RequestedTradingEvent{
					Id:        event.Id.Uint64(),
					Requester: event.Requester.Hex(),
					Product:   event.Product.Uint64(),
				},
			}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		}
	}
}
func (s *tradingContractServiceServer) WatchAcceptedTradingRequestEvent(
	req *WatchAcceptedTradingRequestEventRequest,
	stream TradingContractService_WatchAcceptedTradingRequestEventServer,
) error {
	sink := make(chan *bindings.TradingContractAcceptedTradingRequest)
	defer close(sink)

	sub, err := s.tradingContract.WatchAcceptedTradingRequestEvent(
		&bind.WatchOpts{
			Context: stream.Context(),
		},
		sink,
		contracts.UInt64ToBigInt(req.Ids),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			if err := stream.Send(&WatchAcceptedTradingRequestEventResponse{
				Event: &domain.AcceptedTradingRequestEvent{
					Id:    event.Id.Uint64(),
					Trade: event.Trade.Uint64(),
				},
			}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		}
	}
}
func (s *tradingContractServiceServer) WatchDeclinedTradingRequestEvent(
	req *WatchDeclinedTradingRequestEventRequest,
	stream TradingContractService_WatchDeclinedTradingRequestEventServer,
) error {
	sink := make(chan *bindings.TradingContractDeclinedTradingRequest)
	defer close(sink)

	sub, err := s.tradingContract.WatchDeclinedTradingRequestEvent(
		&bind.WatchOpts{
			Context: stream.Context(),
		},
		sink,
		contracts.UInt64ToBigInt(req.Ids),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			if err := stream.Send(&WatchDeclinedTradingRequestEventResponse{
				Event: &domain.DeclinedTradingRequestEvent{
					Id: event.Id.Uint64(),
				},
			}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		}
	}
}
