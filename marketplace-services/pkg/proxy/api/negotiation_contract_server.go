package api

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/contracts/bindings"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/services"
	"math/big"
)

type negotiationContractServiceServer struct {
	UnimplementedNegotiationContractServiceServer
	negotiationContractService services.NegotiationContractService
	negotiationContract        contracts.NegotiationContract
}

func NewNegotiationContractServiceServer(
	negotiationContractService services.NegotiationContractService,
	negotiationContract contracts.NegotiationContract,
) *negotiationContractServiceServer {
	return &negotiationContractServiceServer{
		negotiationContractService: negotiationContractService,
		negotiationContract:        negotiationContract,
	}
}

func (s *negotiationContractServiceServer) RequestNegotiation(
	ctx context.Context,
	req *RequestNegotiationRequest,
) (*RequestNegotiationResponse, error) {
	tx, err := s.negotiationContractService.RequestNegotiation(ctx, big.NewInt(int64(req.Product)))
	if err != nil {
		return &RequestNegotiationResponse{}, err
	}
	return &RequestNegotiationResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *negotiationContractServiceServer) AcceptNegotiationRequest(
	ctx context.Context,
	req *AcceptNegotiationRequestRequest,
) (*AcceptNegotiationRequestResponse, error) {
	tx, err := s.negotiationContractService.RequestNegotiation(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &AcceptNegotiationRequestResponse{}, err
	}
	return &AcceptNegotiationRequestResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *negotiationContractServiceServer) DeclineNegotiationRequest(
	ctx context.Context,
	req *DeclineNegotiationRequestRequest,
) (*DeclineNegotiationRequestResponse, error) {
	tx, err := s.negotiationContractService.RequestNegotiation(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &DeclineNegotiationRequestResponse{}, err
	}
	return &DeclineNegotiationRequestResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *negotiationContractServiceServer) FindNegotiationRequestByIndex(
	ctx context.Context,
	req *FindNegotiationRequestByIndexRequest,
) (*FindNegotiationRequestByIndexResponse, error) {
	r, err := s.negotiationContractService.FindNegotiationRequestByIndex(ctx, big.NewInt(int64(req.Index)))
	if err != nil {
		return &FindNegotiationRequestByIndexResponse{}, err
	}
	return &FindNegotiationRequestByIndexResponse{Request: NegotiationRequestToGrpcNegotiationRequest(r)}, err
}
func (s *negotiationContractServiceServer) FindNegotiationByIndex(
	ctx context.Context,
	req *FindNegotiationByIndexRequest,
) (*FindNegotiationByIndexResponse, error) {
	n, err := s.negotiationContractService.FindNegotiationByIndex(ctx, big.NewInt(int64(req.Index)))
	if err != nil {
		return &FindNegotiationByIndexResponse{}, err
	}
	return &FindNegotiationByIndexResponse{Negotiation: NegotiationToGrpcNegotiation(n)}, err
}
func (s *negotiationContractServiceServer) FindNegotiationRequestById(
	ctx context.Context,
	req *FindNegotiationRequestByIdRequest,
) (*FindNegotiationRequestByIdResponse, error) {
	r, err := s.negotiationContractService.FindNegotiationRequestById(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &FindNegotiationRequestByIdResponse{}, err
	}
	return &FindNegotiationRequestByIdResponse{Request: NegotiationRequestToGrpcNegotiationRequest(r)}, err
}
func (s *negotiationContractServiceServer) FindNegotiationById(
	ctx context.Context,
	req *FindNegotiationByIdRequest,
) (*FindNegotiationByIdResponse, error) {
	n, err := s.negotiationContractService.FindNegotiationById(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &FindNegotiationByIdResponse{}, err
	}
	return &FindNegotiationByIdResponse{Negotiation: NegotiationToGrpcNegotiation(n)}, err
}
func (s *negotiationContractServiceServer) CountNegotiationRequests(
	ctx context.Context,
	_ *CountNegotiationRequestsRequest,
) (*CountNegotiationRequestsResponse, error) {
	c, err := s.negotiationContractService.CountNegotiationRequests(ctx)
	if err != nil {
		return &CountNegotiationRequestsResponse{}, err
	}
	return &CountNegotiationRequestsResponse{Counter: c.Uint64()}, err
}
func (s *negotiationContractServiceServer) CountNegotiations(
	ctx context.Context,
	_ *CountNegotiationsRequest,
) (*CountNegotiationsResponse, error) {
	c, err := s.negotiationContractService.CountNegotiations(ctx)
	if err != nil {
		return &CountNegotiationsResponse{}, err
	}
	return &CountNegotiationsResponse{Counter: c.Uint64()}, err
}
func (s *negotiationContractServiceServer) ExistsNegotiationRequestById(
	ctx context.Context,
	req *ExistsNegotiationRequestByIdRequest,
) (*ExistsNegotiationRequestByIdResponse, error) {
	exists, err := s.negotiationContractService.ExistsNegotiationRequestById(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &ExistsNegotiationRequestByIdResponse{}, err
	}
	return &ExistsNegotiationRequestByIdResponse{Exists: exists}, err
}
func (s *negotiationContractServiceServer) ExistsNegotiationById(
	ctx context.Context,
	req *ExistsNegotiationByIdRequest,
) (*ExistsNegotiationByIdResponse, error) {
	exists, err := s.negotiationContractService.ExistsNegotiationById(ctx, big.NewInt(int64(req.Id)))
	if err != nil {
		return &ExistsNegotiationByIdResponse{}, err
	}
	return &ExistsNegotiationByIdResponse{Exists: exists}, err
}
func (s *negotiationContractServiceServer) WatchRequestedNegotiationEvent(
	req *WatchRequestedNegotiationEventRequest,
	stream NegotiationContractService_WatchRequestedNegotiationEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.NegotiationContractRequestedNegotiation)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.negotiationContract.WatchRequestedNegotiationEvent(
		watchOpts,
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
		case e := <-sink:
			event := &domain.RequestedNegotiationEvent{
				Requester: e.Requester.Hex(),
				Id:        e.Id.Uint64(),
				Product:   e.Product.Uint64(),
			}
			if err := stream.Send(&WatchRequestedNegotiationEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *negotiationContractServiceServer) WatchAcceptedNegotiationRequestEvent(
	req *WatchAcceptedNegotiationRequestEventRequest,
	stream NegotiationContractService_WatchAcceptedNegotiationRequestEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.NegotiationContractAcceptedNegotiationRequest)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.negotiationContract.WatchAcceptedNegotiationRequestEvent(
		watchOpts,
		sink,
		contracts.UInt64ToBigInt(req.Ids),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.AcceptedNegotiationRequestEvent{
				Id:          e.Id.Uint64(),
				Negotiation: e.Negotiation.Uint64(),
			}
			if err := stream.Send(&WatchAcceptedNegotiationRequestEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *negotiationContractServiceServer) WatchDeclinedNegotiationRequestEvent(
	req *WatchDeclinedNegotiationRequestEventRequest,
	stream NegotiationContractService_WatchDeclinedNegotiationRequestEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.NegotiationContractDeclinedNegotiationRequest)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.negotiationContract.WatchDeclinedNegotiationRequestEvent(
		watchOpts,
		sink,
		contracts.UInt64ToBigInt(req.Ids),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.DeclinedNegotiationRequestEvent{
				Id: e.Id.Uint64(),
			}
			if err := stream.Send(&WatchDeclinedNegotiationRequestEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
