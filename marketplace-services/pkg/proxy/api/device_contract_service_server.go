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

type deviceContractServiceServer struct {
	UnimplementedDeviceContractServiceServer
	deviceContractService services.DeviceContractService
	deviceContract        contracts.DeviceContract
}

func NewDeviceContractServiceServer(
	deviceContractService services.DeviceContractService,
	deviceContract contracts.DeviceContract,
) *deviceContractServiceServer {
	return &deviceContractServiceServer{deviceContractService: deviceContractService, deviceContract: deviceContract}
}

func (s *deviceContractServiceServer) CreateDevice(
	ctx context.Context,
	req *CreateDeviceRequest,
) (*CreateDeviceResponse, error) {
	tx, err := s.deviceContractService.CreateDevice(ctx, DeviceFromGrpcDevice(req.Device))
	if err != nil {
		return &CreateDeviceResponse{}, err
	}
	return &CreateDeviceResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *deviceContractServiceServer) UpdateDevice(
	ctx context.Context,
	req *UpdateDeviceRequest,
) (*UpdateDeviceResponse, error) {
	tx, err := s.deviceContractService.UpdateDevice(ctx, DeviceFromGrpcDevice(req.Device))
	if err != nil {
		return &UpdateDeviceResponse{}, err
	}
	return &UpdateDeviceResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *deviceContractServiceServer) RemoveDevice(
	ctx context.Context,
	req *RemoveDeviceRequest,
) (*RemoveDeviceResponse, error) {
	tx, err := s.deviceContractService.RemoveDevice(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &RemoveDeviceResponse{}, err
	}
	return &RemoveDeviceResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *deviceContractServiceServer) FindDeviceByIndex(
	ctx context.Context,
	req *FindDeviceByIndexRequest,
) (*FindDeviceByIndexResponse, error) {
	d, err := s.deviceContractService.FindDeviceByIndex(ctx, big.NewInt(req.Index))
	if err != nil {
		return &FindDeviceByIndexResponse{}, err
	}
	return &FindDeviceByIndexResponse{Device: DeviceToGrpcDevice(d)}, err
}
func (s *deviceContractServiceServer) FindDeviceByAddress(
	ctx context.Context,
	req *FindDeviceByAddressRequest,
) (*FindDeviceByAddressResponse, error) {
	d, err := s.deviceContractService.FindDeviceByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindDeviceByAddressResponse{}, err
	}
	return &FindDeviceByAddressResponse{Device: DeviceToGrpcDevice(d)}, err
}
func (s *deviceContractServiceServer) FindProductsOfDeviceByAddress(
	ctx context.Context,
	req *FindProductsOfDeviceByAddressRequest,
) (*FindProductsOfDeviceByAddressResponse, error) {
	p, err := s.deviceContractService.FindProductsOfDeviceByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindProductsOfDeviceByAddressResponse{}, err
	}
	return &FindProductsOfDeviceByAddressResponse{Products: contracts.BigIntToInt64(p)}, err
}
func (s *deviceContractServiceServer) FindNegotiationRequestsOfDeviceByAddress(
	ctx context.Context,
	req *FindNegotiationRequestsOfDeviceByAddressRequest,
) (*FindNegotiationRequestsOfDeviceByAddressResponse, error) {
	nr, err := s.deviceContractService.FindNegotiationRequestsOfDeviceByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindNegotiationRequestsOfDeviceByAddressResponse{}, err
	}
	return &FindNegotiationRequestsOfDeviceByAddressResponse{NegotiationRequests: contracts.BigIntToInt64(nr)}, err
}
func (s *deviceContractServiceServer) FindTradingRequestsOfDeviceByAddress(
	ctx context.Context,
	req *FindTradingRequestsOfDeviceByAddressRequest,
) (*FindTradingRequestsOfDeviceByAddressResponse, error) {
	tr, err := s.deviceContractService.FindTradingRequestsOfDeviceByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindTradingRequestsOfDeviceByAddressResponse{}, err
	}
	return &FindTradingRequestsOfDeviceByAddressResponse{TradingRequests: contracts.BigIntToInt64(tr)}, err
}
func (s *deviceContractServiceServer) FindNegotiationsOfDeviceByAddress(
	ctx context.Context,
	req *FindNegotiationsOfDeviceByAddressRequest,
) (*FindNegotiationsOfDeviceByAddressResponse, error) {
	n, err := s.deviceContractService.FindNegotiationsOfDeviceByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindNegotiationsOfDeviceByAddressResponse{}, err
	}
	return &FindNegotiationsOfDeviceByAddressResponse{Negotiations: contracts.BigIntToInt64(n)}, err
}
func (s *deviceContractServiceServer) FindTradesOfDeviceByAddress(
	ctx context.Context,
	req *FindTradesOfDeviceByAddressRequest,
) (*FindTradesOfDeviceByAddressResponse, error) {
	t, err := s.deviceContractService.FindTradesOfDeviceByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindTradesOfDeviceByAddressResponse{}, err
	}
	return &FindTradesOfDeviceByAddressResponse{Trades: contracts.BigIntToInt64(t)}, err
}
func (s *deviceContractServiceServer) IsDeviceOwnedByUser(
	ctx context.Context,
	req *IsDeviceOwnedByUserRequest,
) (*IsDeviceOwnedByUserResponse, error) {
	owned, err := s.deviceContractService.IsDeviceOwnedByUser(ctx, common.HexToAddress(req.Address), common.HexToAddress(req.User))
	if err != nil {
		return &IsDeviceOwnedByUserResponse{}, err
	}
	return &IsDeviceOwnedByUserResponse{Owned: owned}, err
}
func (s *deviceContractServiceServer) CountDevices(
	ctx context.Context,
	_ *CountDevicesRequest,
) (*CountDevicesResponse, error) {
	c, err := s.deviceContractService.CountDevices(ctx)
	if err != nil {
		return &CountDevicesResponse{}, err
	}
	return &CountDevicesResponse{Counter: c.Int64()}, err
}
func (s *deviceContractServiceServer) ExistsDeviceByAddress(
	ctx context.Context,
	req *ExistsDeviceByAddressRequest,
) (*ExistsDeviceByAddressResponse, error) {
	exists, err := s.deviceContractService.ExistsDeviceByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &ExistsDeviceByAddressResponse{}, err
	}
	return &ExistsDeviceByAddressResponse{Exists: exists}, err
}
func (s *deviceContractServiceServer) ExistsDeviceByAddressAndDeleted(
	ctx context.Context,
	req *ExistsDeviceByAddressAndDeletedRequest,
) (*ExistsDeviceByAddressAndDeletedResponse, error) {
	exists, err := s.deviceContractService.ExistsDeviceByAddressAndDeleted(ctx, common.HexToAddress(req.Address), req.Deleted)
	if err != nil {
		return &ExistsDeviceByAddressAndDeletedResponse{}, err
	}
	return &ExistsDeviceByAddressAndDeletedResponse{Exists: exists}, err
}
func (s *deviceContractServiceServer) WatchCreatedDeviceEvent(
	req *WatchCreatedDeviceEventRequest,
	stream DeviceContractService_WatchCreatedDeviceEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.DeviceContractCreatedDevice)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.deviceContract.WatchCreatedDeviceEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.CreatedDeviceEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchCreatedDeviceEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *deviceContractServiceServer) WatchUpdatedDeviceEvent(
	req *WatchUpdatedDeviceEventRequest,
	stream DeviceContractService_WatchUpdatedDeviceEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.DeviceContractUpdatedDevice)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.deviceContract.WatchUpdatedDeviceEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.UpdatedDeviceEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchUpdatedDeviceEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *deviceContractServiceServer) WatchRemovedDeviceEvent(
	req *WatchRemovedDeviceEventRequest,
	stream DeviceContractService_WatchRemovedDeviceEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.DeviceContractRemovedDevice)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.deviceContract.WatchRemovedDeviceEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.RemovedDeviceEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchRemovedDeviceEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
