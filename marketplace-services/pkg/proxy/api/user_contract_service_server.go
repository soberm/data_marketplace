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

type userContractServiceServer struct {
	UnimplementedUserContractServiceServer
	userContractService services.UserContractService
	userContract        contracts.UserContract
}

func NewUserContractServiceServer(
	userContractService services.UserContractService,
	userContract contracts.UserContract,
) *userContractServiceServer {
	return &userContractServiceServer{userContractService: userContractService, userContract: userContract}
}

func (s *userContractServiceServer) CreateUser(
	ctx context.Context,
	req *CreateUserRequest,
) (*CreateUserResponse, error) {
	tx, err := s.userContractService.CreateUser(ctx, UserFromGrpcUser(req.User))
	if err != nil {
		return &CreateUserResponse{}, err
	}
	return &CreateUserResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}

func (s *userContractServiceServer) UpdateUser(
	ctx context.Context,
	req *UpdateUserRequest,
) (*UpdateUserResponse, error) {
	tx, err := s.userContractService.UpdateUser(ctx, UserFromGrpcUser(req.User))
	if err != nil {
		return &UpdateUserResponse{}, err
	}
	return &UpdateUserResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}

func (s *userContractServiceServer) RemoveUser(
	ctx context.Context,
	_ *RemoveUserRequest,
) (*RemoveUserResponse, error) {
	tx, err := s.userContractService.RemoveUser(ctx)
	if err != nil {
		return &RemoveUserResponse{}, err
	}
	return &RemoveUserResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}

func (s *userContractServiceServer) FindUserByIndex(
	ctx context.Context,
	req *FindUserByIndexRequest,
) (*FindUserByIndexResponse, error) {
	u, err := s.userContractService.FindUserByIndex(ctx, big.NewInt(req.Index))
	if err != nil {
		return &FindUserByIndexResponse{}, err
	}
	return &FindUserByIndexResponse{User: UserToGrpcUser(u)}, err
}

func (s *userContractServiceServer) FindUserByAddress(
	ctx context.Context,
	req *FindUserByAddressRequest,
) (*FindUserByAddressResponse, error) {
	u, err := s.userContractService.FindUserByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &FindUserByAddressResponse{}, err
	}
	return &FindUserByAddressResponse{User: UserToGrpcUser(u)}, err
}

func (s *userContractServiceServer) CountUsers(
	ctx context.Context,
	_ *CountUsersRequest,
) (*CountUsersResponse, error) {
	counter, err := s.userContractService.CountUsers(ctx)
	if err != nil {
		return &CountUsersResponse{}, err
	}
	return &CountUsersResponse{Counter: counter.Int64()}, err
}

func (s *userContractServiceServer) ExistsUserByAddress(
	ctx context.Context,
	req *ExistsUserByAddressRequest,
) (*ExistsUserByAddressResponse, error) {
	exists, err := s.userContractService.ExistsUserByAddress(ctx, common.HexToAddress(req.Address))
	if err != nil {
		return &ExistsUserByAddressResponse{}, err
	}
	return &ExistsUserByAddressResponse{Exists: exists}, err
}

func (s *userContractServiceServer) ExistsUserByAddressAndDeleted(
	ctx context.Context,
	req *ExistsUserByAddressAndDeletedRequest,
) (*ExistsUserByAddressAndDeletedResponse, error) {
	exists, err := s.userContractService.ExistsUserByAddressAndDeleted(ctx, common.HexToAddress(req.Address), req.Deleted)
	if err != nil {
		return &ExistsUserByAddressAndDeletedResponse{}, err
	}
	return &ExistsUserByAddressAndDeletedResponse{Exists: exists}, err
}

func (s *userContractServiceServer) WatchCreatedUserEvent(
	req *WatchCreatedUserEventRequest,
	stream UserContractService_WatchCreatedUserEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.UserContractCreatedUser)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.userContract.WatchCreatedUserEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.CreatedUserEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchCreatedUserEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}

func (s *userContractServiceServer) WatchUpdatedUserEvent(
	req *WatchUpdatedUserEventRequest,
	stream UserContractService_WatchUpdatedUserEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.UserContractUpdatedUser)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.userContract.WatchUpdatedUserEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.UpdatedUserEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchUpdatedUserEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}

func (s *userContractServiceServer) WatchRemovedUserEvent(
	req *WatchRemovedUserEventRequest,
	stream UserContractService_WatchRemovedUserEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.UserContractRemovedUser)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.userContract.WatchRemovedUserEvent(watchOpts, sink, contracts.HexToAddresses(req.Addresses))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.RemovedUserEvent{Address: e.Addr.Hex()}
			if err := stream.Send(&WatchRemovedUserEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
