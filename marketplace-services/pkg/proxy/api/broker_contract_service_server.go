package api

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/contracts/bindings"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/services"
	"math/big"
)

type brokerContractServiceServer struct {
	UnimplementedBrokerContractServiceServer
	brokerContractService services.BrokerContractService
	brokerContract        contracts.BrokerContract
}

func NewBrokerContractServiceServer(
	brokerContractService services.BrokerContractService,
	brokerContract contracts.BrokerContract,
) *brokerContractServiceServer {
	return &brokerContractServiceServer{brokerContractService: brokerContractService, brokerContract: brokerContract}
}

func (s *brokerContractServiceServer) CreateBroker(
	ctx context.Context,
	req *CreateBrokerRequest,
) (*CreateBrokerResponse, error) {
	tx, err := s.brokerContractService.CreateBroker(ctx, BrokerFromGrpcBroker(req.Broker))
	if err != nil {
		return &CreateBrokerResponse{}, err
	}
	return &CreateBrokerResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *brokerContractServiceServer) UpdateBroker(
	ctx context.Context,
	req *UpdateBrokerRequest,
) (*UpdateBrokerResponse, error) {
	tx, err := s.brokerContractService.UpdateBroker(ctx, BrokerFromGrpcBroker(req.Broker))
	if err != nil {
		return &UpdateBrokerResponse{}, err
	}
	return &UpdateBrokerResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *brokerContractServiceServer) RemoveBroker(
	ctx context.Context,
	req *RemoveBrokerRequest,
) (*RemoveBrokerResponse, error) {
	tx, err := s.brokerContractService.RemoveBroker(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &RemoveBrokerResponse{}, err
	}
	return &RemoveBrokerResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *brokerContractServiceServer) FindBrokerByIndex(
	ctx context.Context,
	req *FindBrokerByIndexRequest,
) (*FindBrokerByIndexResponse, error) {
	b, err := s.brokerContractService.FindBrokerByIndex(ctx, big.NewInt(req.Index))
	if err != nil {
		return &FindBrokerByIndexResponse{}, err
	}
	return &FindBrokerByIndexResponse{Broker: BrokerToGrpcBroker(b)}, err
}
func (s *brokerContractServiceServer) FindBrokerByAddress(
	ctx context.Context,
	req *FindBrokerByAddressRequest,
) (*FindBrokerByAddressResponse, error) {
	b, err := s.brokerContractService.FindBrokerByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindBrokerByAddressResponse{}, err
	}
	return &FindBrokerByAddressResponse{Broker: BrokerToGrpcBroker(b)}, err
}
func (s *brokerContractServiceServer) CountBrokers(
	ctx context.Context,
	_ *CountBrokersRequest,
) (*CountBrokersResponse, error) {
	c, err := s.brokerContractService.CountBrokers(ctx)
	if err != nil {
		return &CountBrokersResponse{}, err
	}
	return &CountBrokersResponse{Counter: c.Int64()}, err
}
func (s *brokerContractServiceServer) ExistsBrokerByAddress(
	ctx context.Context,
	req *ExistsBrokerByAddressRequest,
) (*ExistsBrokerByAddressResponse, error) {
	exists, err := s.brokerContractService.ExistsBrokerByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &ExistsBrokerByAddressResponse{}, err
	}
	return &ExistsBrokerByAddressResponse{Exists: exists}, err
}
func (s *brokerContractServiceServer) ExistsBrokerByAddressAndDeleted(
	ctx context.Context,
	req *ExistsBrokerByAddressAndDeletedRequest,
) (*ExistsBrokerByAddressAndDeletedResponse, error) {
	exists, err := s.brokerContractService.ExistsBrokerByAddressAndDeleted(ctx, common.HexToAddress(req.Address), req.Deleted)
	if err != nil {
		return &ExistsBrokerByAddressAndDeletedResponse{}, err
	}
	return &ExistsBrokerByAddressAndDeletedResponse{Exists: exists}, err
}
func (s *brokerContractServiceServer) WatchCreatedBrokerEvent(
	req *WatchCreatedBrokerEventRequest,
	stream BrokerContractService_WatchCreatedBrokerEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.BrokerContractCreatedBroker)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.brokerContract.WatchCreatedBrokerEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.CreatedBrokerEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchCreatedBrokerEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *brokerContractServiceServer) WatchUpdatedBrokerEvent(
	req *WatchUpdatedBrokerEventRequest,
	stream BrokerContractService_WatchUpdatedBrokerEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.BrokerContractUpdatedBroker)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.brokerContract.WatchUpdatedBrokerEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.UpdatedBrokerEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchUpdatedBrokerEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *brokerContractServiceServer) WatchRemovedBrokerEvent(
	req *WatchRemovedBrokerEventRequest,
	stream BrokerContractService_WatchRemovedBrokerEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.BrokerContractRemovedBroker)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.brokerContract.WatchRemovedBrokerEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.RemovedBrokerEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchRemovedBrokerEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
