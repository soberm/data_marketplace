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

type productContractServiceServer struct {
	UnimplementedProductContractServiceServer
	productContractService services.ProductContractService
	productContract        contracts.ProductContract
}

func NewProductContractServiceServer(
	productContractService services.ProductContractService,
	productContract contracts.ProductContract,
) *productContractServiceServer {
	return &productContractServiceServer{
		productContractService: productContractService,
		productContract:        productContract,
	}
}

func (s *productContractServiceServer) CreateProduct(
	ctx context.Context,
	req *CreateProductRequest,
) (*CreateProductResponse, error) {
	tx, err := s.productContractService.CreateProduct(ctx, ProductFromGrpcProduct(req.Product))
	if err != nil {
		return &CreateProductResponse{}, err
	}
	return &CreateProductResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *productContractServiceServer) UpdateProduct(
	ctx context.Context,
	req *UpdateProductRequest,
) (*UpdateProductResponse, error) {
	tx, err := s.productContractService.UpdateProduct(ctx, ProductFromGrpcProduct(req.Product))
	if err != nil {
		return &UpdateProductResponse{}, err
	}
	return &UpdateProductResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *productContractServiceServer) RemoveProduct(
	ctx context.Context,
	req *RemoveProductRequest,
) (*RemoveProductResponse, error) {
	tx, err := s.productContractService.RemoveProduct(ctx, big.NewInt(req.Id))
	if err != nil {
		return &RemoveProductResponse{}, err
	}
	return &RemoveProductResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}
func (s *productContractServiceServer) FindProductByIndex(
	ctx context.Context,
	req *FindProductByIndexRequest,
) (*FindProductByIndexResponse, error) {
	p, err := s.productContractService.FindProductByIndex(ctx, big.NewInt(req.Index))
	if err != nil {
		return &FindProductByIndexResponse{}, err
	}
	return &FindProductByIndexResponse{Product: ProductToGrpcProduct(p)}, err
}
func (s *productContractServiceServer) FindProductById(
	ctx context.Context,
	req *FindProductByIdRequest,
) (*FindProductByIdResponse, error) {
	p, err := s.productContractService.FindProductByIndex(ctx, big.NewInt(req.Id))
	if err != nil {
		return &FindProductByIdResponse{}, err
	}
	return &FindProductByIdResponse{Product: ProductToGrpcProduct(p)}, err
}
func (s *productContractServiceServer) FindNegotiationRequestsOfProductById(
	ctx context.Context,
	req *FindNegotiationRequestsOfProductByIdRequest,
) (*FindNegotiationRequestsOfProductByIdResponse, error) {
	nr, err := s.productContractService.FindNegotiationRequestsOfProductById(ctx, big.NewInt(req.Id))
	if err != nil {
		return &FindNegotiationRequestsOfProductByIdResponse{}, err
	}
	return &FindNegotiationRequestsOfProductByIdResponse{NegotiationRequests: contracts.BigIntToInt64(nr)}, err
}
func (s *productContractServiceServer) FindTradingRequestsOfProductById(
	ctx context.Context,
	req *FindTradingRequestsOfProductByIdRequest,
) (*FindTradingRequestsOfProductByIdResponse, error) {
	tr, err := s.productContractService.FindTradingRequestsOfProductById(ctx, big.NewInt(req.Id))
	if err != nil {
		return &FindTradingRequestsOfProductByIdResponse{}, err
	}
	return &FindTradingRequestsOfProductByIdResponse{TradingRequests: contracts.BigIntToInt64(tr)}, err
}
func (s *productContractServiceServer) FindNegotiationsOfProductById(
	ctx context.Context,
	req *FindNegotiationsOfProductByIdRequest,
) (*FindNegotiationsOfProductByIdResponse, error) {
	n, err := s.productContractService.FindNegotiationsOfProductById(ctx, big.NewInt(req.Id))
	if err != nil {
		return &FindNegotiationsOfProductByIdResponse{}, err
	}
	return &FindNegotiationsOfProductByIdResponse{Negotiations: contracts.BigIntToInt64(n)}, err
}
func (s *productContractServiceServer) FindTradesOfProductById(
	ctx context.Context,
	req *FindTradesOfProductByIdRequest,
) (*FindTradesOfProductByIdResponse, error) {
	t, err := s.productContractService.FindTradesOfProductById(ctx, big.NewInt(req.Id))
	if err != nil {
		return &FindTradesOfProductByIdResponse{}, err
	}
	return &FindTradesOfProductByIdResponse{Trades: contracts.BigIntToInt64(t)}, err
}
func (s *productContractServiceServer) FindCostOfProductById(
	ctx context.Context,
	req *FindCostOfProductByIdRequest,
) (*FindCostOfProductByIdResponse, error) {
	c, err := s.productContractService.FindCostOfProductById(ctx, big.NewInt(req.Id))
	if err != nil {
		return &FindCostOfProductByIdResponse{}, err
	}
	return &FindCostOfProductByIdResponse{Cost: c.Int64()}, err
}
func (s *productContractServiceServer) IsProductOwnedByDevice(
	ctx context.Context,
	req *IsProductOwnedByDeviceRequest,
) (*IsProductOwnedByDeviceResponse, error) {
	owned, err := s.productContractService.IsProductOwnedByDevice(ctx, big.NewInt(req.Id), common.HexToAddress(req.Device))
	if err != nil {
		return &IsProductOwnedByDeviceResponse{}, err
	}
	return &IsProductOwnedByDeviceResponse{Owned: owned}, err
}
func (s *productContractServiceServer) CountProducts(
	ctx context.Context,
	_ *CountProductsRequest,
) (*CountProductsResponse, error) {
	c, err := s.productContractService.CountProducts(ctx)
	if err != nil {
		return &CountProductsResponse{}, err
	}
	return &CountProductsResponse{Counter: c.Int64()}, err
}
func (s *productContractServiceServer) ExistsProductById(
	ctx context.Context,
	req *ExistsProductByIdRequest,
) (*ExistsProductByIdResponse, error) {
	exists, err := s.productContractService.ExistsProductById(ctx, big.NewInt(req.Id))
	if err != nil {
		return &ExistsProductByIdResponse{}, err
	}
	return &ExistsProductByIdResponse{Exists: exists}, err
}
func (s *productContractServiceServer) ExistsProductByIdAndDeleted(
	ctx context.Context,
	req *ExistsProductByIdAndDeletedRequest,
) (*ExistsProductByIdAndDeletedResponse, error) {
	exists, err := s.productContractService.ExistsProductByIdAndDeleted(ctx, big.NewInt(req.Id), req.Deleted)
	if err != nil {
		return &ExistsProductByIdAndDeletedResponse{}, err
	}
	return &ExistsProductByIdAndDeletedResponse{Exists: exists}, err
}
func (s *productContractServiceServer) WatchCreatedProductEvent(
	req *WatchCreatedProductEventRequest,
	stream ProductContractService_WatchCreatedProductEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.ProductContractCreatedProduct)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.productContract.WatchCreatedProductEvent(watchOpts, sink, contracts.HexToAddresses(req.Users))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.CreatedProductEvent{Id: e.Id.Uint64(), User: e.User.Hex()}
			if err := stream.Send(&WatchCreatedProductEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *productContractServiceServer) WatchUpdatedProductEvent(
	req *WatchUpdatedProductEventRequest,
	stream ProductContractService_WatchUpdatedProductEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.ProductContractUpdatedProduct)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.productContract.WatchUpdatedProductEvent(
		watchOpts,
		sink,
		contracts.UInt64ToBigInt(req.Ids),
		contracts.HexToAddresses(req.Users),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.UpdatedProductEvent{Id: e.Id.Uint64(), User: e.User.Hex()}
			if err := stream.Send(&WatchUpdatedProductEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
func (s *productContractServiceServer) WatchRemovedProductEvent(
	req *WatchRemovedProductEventRequest,
	stream ProductContractService_WatchRemovedProductEventServer,
) error {
	ctx := stream.Context()
	sink := make(chan *bindings.ProductContractRemovedProduct)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.productContract.WatchRemovedProductEvent(
		watchOpts,
		sink,
		contracts.UInt64ToBigInt(req.Ids),
		contracts.HexToAddresses(req.Users),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.RemovedProductEvent{Id: e.Id.Uint64(), User: e.User.Hex()}
			if err := stream.Send(&WatchRemovedProductEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
