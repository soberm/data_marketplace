package api

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/contracts/bindings"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/services"
	"math/big"
)

type settlementContractServiceServer struct {
	UnimplementedSettlementContractServiceServer
	logger        logrus.FieldLogger
	walletService services.WalletService
	keyStore      *keystore.KeyStore
	ethClient     *ethclient.Client
}

func NewSettlementContractServiceServer(
	logger logrus.FieldLogger,
	walletService services.WalletService,
	keyStore *keystore.KeyStore,
	ethClient *ethclient.Client,
) *settlementContractServiceServer {
	return &settlementContractServiceServer{
		logger:        logger,
		walletService: walletService,
		keyStore:      keyStore,
		ethClient:     ethClient,
	}
}

func (s *settlementContractServiceServer) Deposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &DepositResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	tx, err := service.Deposit(ctx, big.NewInt(int64(req.Value)))
	if err != nil {
		return &DepositResponse{}, err
	}
	return &DepositResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}

func (s *settlementContractServiceServer) SettleTrade(ctx context.Context, req *SettleTradeRequest) (*SettleTradeResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &SettleTradeResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	tx, err := service.SettleTrade(ctx, big.NewInt(int64(req.Counter)))
	if err != nil {
		return &SettleTradeResponse{}, err
	}
	return &SettleTradeResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}

func (s *settlementContractServiceServer) ResolveDispute(ctx context.Context, req *ResolveDisputeRequest) (*ResolveDisputeResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &ResolveDisputeResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	tx, err := service.ResolveDispute(ctx, big.NewInt(int64(req.Counter)))
	if err != nil {
		return &ResolveDisputeResponse{}, err
	}
	return &ResolveDisputeResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}

func (s *settlementContractServiceServer) ResolveTimeout(ctx context.Context, req *ResolveTimeoutRequest) (*ResolveTimeoutResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &ResolveTimeoutResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	tx, err := service.ResolveTimeout(ctx)
	if err != nil {
		return &ResolveTimeoutResponse{}, err
	}
	return &ResolveTimeoutResponse{Transaction: TransactionToGrpcTransaction(tx)}, err
}

func (s *settlementContractServiceServer) GetProviderCounter(ctx context.Context, req *GetProviderCounterRequest) (*GetProviderCounterResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &GetProviderCounterResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	counter, err := service.GetProviderCounter(ctx)
	if err != nil {
		return &GetProviderCounterResponse{}, err
	}
	return &GetProviderCounterResponse{Counter: CounterToGrpcCounter(counter)}, err
}

func (s *settlementContractServiceServer) GetConsumerCounter(ctx context.Context, req *GetConsumerCounterRequest) (*GetConsumerCounterResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &GetConsumerCounterResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	counter, err := service.GetConsumerCounter(ctx)
	if err != nil {
		return &GetConsumerCounterResponse{}, err
	}
	return &GetConsumerCounterResponse{Counter: CounterToGrpcCounter(counter)}, err
}

func (s *settlementContractServiceServer) GetBrokerCounter(ctx context.Context, req *GetBrokerCounterRequest) (*GetBrokerCounterResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &GetBrokerCounterResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	counter, err := service.GetBrokerCounter(ctx)
	if err != nil {
		return &GetBrokerCounterResponse{}, err
	}
	return &GetBrokerCounterResponse{Counter: CounterToGrpcCounter(counter)}, err
}

func (s *settlementContractServiceServer) GetSettlement(ctx context.Context, req *GetSettlementRequest) (*GetSettlementResponse, error) {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return &GetSettlementResponse{}, err
	}
	service := services.NewSettlementContractServiceImpl(s.logger, s.walletService, s.keyStore, contract)
	settlement, err := service.GetSettlement(ctx)
	if err != nil {
		return &GetSettlementResponse{}, err
	}
	return &GetSettlementResponse{Settlement: SettlementToGrpcSettlement(settlement)}, err
}

func (s *settlementContractServiceServer) WatchDepositedEvent(
	req *WatchDepositedEventRequest,
	stream SettlementContractService_WatchDepositedEventServer,
) error {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return err
	}

	ctx := stream.Context()
	sink := make(chan *bindings.SettlementContractDeposited)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := contract.WatchDepositedEvent(watchOpts, sink, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.DepositedEvent{Payee: e.Depositor.Hex(), Amount: e.Amount.Uint64()}
			if err := stream.Send(&WatchDepositedEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}

func (s *settlementContractServiceServer) WatchSettledEvent(req *WatchSettledEventRequest, stream SettlementContractService_WatchSettledEventServer) error {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return err
	}

	ctx := stream.Context()
	sink := make(chan *bindings.SettlementContractSettled)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := contract.WatchSettledEvent(watchOpts, sink)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.SettledEvent{
				ActualCost: e.ActualCost.Uint64(),
				Provider:   e.Provider.Uint64(),
				Consumer:   e.Consumer.Uint64(),
				Broker:     e.Broker.Uint64(),
			}
			if err := stream.Send(&WatchSettledEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}

func (s *settlementContractServiceServer) WatchDisputeEvent(req *WatchDisputeEventRequest, stream SettlementContractService_WatchDisputeEventServer) error {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return err
	}

	ctx := stream.Context()
	sink := make(chan *bindings.SettlementContractDispute)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := contract.WatchDisputeEvent(watchOpts, sink)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.DisputeEvent{
				ProviderCounter: e.ProviderCounter.Uint64(),
				ConsumerCounter: e.ConsumerCounter.Uint64(),
			}
			if err := stream.Send(&WatchDisputeEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}

func (s *settlementContractServiceServer) WatchCounterSetEvent(req *WatchCounterSetEventRequest, stream SettlementContractService_WatchCounterSetEventServer) error {
	contract, err := contracts.NewSettlementContractImpl(common.HexToAddress(req.ContractAddress), s.ethClient)
	if err != nil {
		return err
	}

	ctx := stream.Context()
	sink := make(chan *bindings.SettlementContractCounterSet)
	defer close(sink)

	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := contract.WatchCounterSetEvent(watchOpts, sink, contracts.HexToAddresses(req.Setter))
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-sink:
			event := &domain.CounterSetEvent{Setter: e.Setter.Hex(), Counter: e.Counter.Uint64()}
			if err := stream.Send(&WatchCounterSetEventResponse{Event: event}); err != nil {
				return err
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return status.Errorf(codes.Canceled, "%s", ctx.Err())
		}
	}
}
