package broker

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"marketplace-services/pkg/broker/api"
	"marketplace-services/pkg/broker/services"
	"marketplace-services/pkg/contracts"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type broker struct {
	opts           options
	grpcServer     *grpc.Server
	ethClient      *ethclient.Client
	disputeService services.DisputeService
	logger         logrus.FieldLogger
	running        bool
	quit           chan bool
}

func New(opt ...Option) (*broker, error) {
	opts := DefaultOptions()
	for _, o := range opt {
		o.apply(&opts)
	}
	if opts.ConfigFile != "" {
		err := opts.loadConfiguration()
		if err != nil {
			return nil, fmt.Errorf("load configuration: %w", err)
		}
	}

	logger := initLogger(opts)

	ethClient, err := ethclient.Dial(opts.EthConfig.ClientURL)
	if err != nil {
		return nil, fmt.Errorf("dial eth client %s: %w", opts.EthConfig.ClientURL, err)
	}

	ks := keystore.NewKeyStore(
		opts.EthConfig.KeyDir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	productContract, err := contracts.NewProductContractImpl(
		common.HexToAddress(opts.ContractsConfig.ProductContractAddress),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"new product contract with address %s: %w",
			opts.ContractsConfig.ProductContractAddress,
			err,
		)
	}

	tradingContract, err := contracts.NewTradingContractImpl(
		common.HexToAddress(opts.ContractsConfig.TradingContractAddress),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"new trading contract with address %s: %w",
			opts.ContractsConfig.TradingContractAddress,
			err,
		)
	}

	messageService := services.NewMessageServiceImpl(logger)
	messageServiceServer := api.NewMessageServiceServer(messageService)

	discoveryService := services.NewDiscoveryServiceImpl(logger, productContract, ks)
	discoveryServiceServer := api.NewDiscoveryServiceServer(discoveryService)

	disputeService := services.NewDisputeServiceImpl(
		logger,
		ks,
		ethClient,
		tradingContract,
		messageService,
		opts.EthConfig.Account,
		opts.EthConfig.Passphrase,
	)

	grpcServer := initGrpcServer(logger)
	api.RegisterMessageServiceServer(grpcServer, messageServiceServer)
	api.RegisterDiscoveryServiceServer(grpcServer, discoveryServiceServer)

	b := &broker{
		opts:           opts,
		grpcServer:     grpcServer,
		ethClient:      ethClient,
		disputeService: disputeService,
		logger:         logger,
		quit:           make(chan bool, 1),
	}

	return b, nil
}

func initLogger(opts options) logrus.FieldLogger {
	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.Level(opts.LoggingConfig.Verbosity),
	}
	return logger
}

func initGrpcServer(logger logrus.FieldLogger) *grpc.Server {
	entry := logrus.NewEntry(logger.(*logrus.Logger))
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(
				entry,
				grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
			),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(
				entry,
				grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
			),
		)),

	)
	grpc_logrus.ReplaceGrpcLogger(entry)
	return server
}

func (b *broker) Run() error {
	addr := b.opts.Host + ":" + strconv.Itoa(b.opts.Port)

	b.logger.Infof("Starting broker %s", version)

	if len(build) != 0 {
		b.logger.Infof("Built on %s", build)
	}

	if len(gitCommit) != 0 {
		b.logger.Infof("Git commit %s", gitCommit)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("run broker: %w", err)
	}

	b.logger.Infof("Broker listening on %s", lis.Addr())

	b.receiveSignals()

	go func() {
		err := b.disputeService.ResolveDisputes(context.TODO())
		if err != nil {
			b.logger.Errorf("resolve disputes: %v", err)
		}
	}()

	b.running = true
	return b.grpcServer.Serve(lis)
}

func (b *broker) receiveSignals() {
	if b.opts.NoSig {
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go b.handleSignals(sigs)
}

func (b *broker) handleSignals(sigs chan os.Signal) {
	select {
	case <-sigs:
		b.GraceFulStop()
		os.Exit(0)
	case <-b.quit:
		return
	}
}

func (b *broker) GraceFulStop() {
	b.logger.Infof("Initiating graceful stop of broker")
	if !b.running {
		return
	}
	close(b.quit)
	b.grpcServer.GracefulStop()
	b.ethClient.Close()
	b.running = false
	b.logger.Infof("Gracefully stopped broker")
}
