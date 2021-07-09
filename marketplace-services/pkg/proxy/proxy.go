package proxy

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_middleware_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/proxy/api"
	"marketplace-services/pkg/proxy/model"
	"marketplace-services/pkg/proxy/services"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type proxy struct {
	opts   options
	db     *gorm.DB
	logger logrus.FieldLogger

	grpcServer *grpc.Server
	ethClient  *ethclient.Client

	running bool
	quit    chan bool
}

func New(opt ...Option) (*proxy, error) {
	opts := defaultOptions()
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
	db, err := initDb(logger, opts)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	ethClient, err := ethclient.Dial(opts.EthConfig.ClientURL)
	if err != nil {
		return nil, fmt.Errorf("dial eth client %s: %w", opts.EthConfig.ClientURL, err)
	}

	ks := keystore.NewKeyStore(
		opts.EthConfig.KeyDir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	userContract, err := contracts.NewUserContractImpl(
		common.HexToAddress(opts.ContractsConfig.UserContractAddress),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"new user contract with address %s: %w",
			opts.ContractsConfig.UserContractAddress,
			err,
		)
	}

	deviceContract, err := contracts.NewDeviceContractImpl(
		common.HexToAddress(opts.ContractsConfig.DeviceContractAddress),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"new device contract with address %s: %w",
			opts.ContractsConfig.DeviceContractAddress,
			err,
		)
	}

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

	brokerContract, err := contracts.NewBrokerContractImpl(
		common.HexToAddress(opts.ContractsConfig.BrokerContractAddress),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"new broker contract with address %s: %w",
			opts.ContractsConfig.BrokerContractAddress,
			err,
		)
	}

	negotiationContract, err := contracts.NewNegotiationContractImpl(
		common.HexToAddress(opts.ContractsConfig.NegotiationContractAddress),
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"new negotiation contract with address %s: %w",
			opts.ContractsConfig.NegotiationContractAddress,
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

	accountService := services.NewAccountServiceImpl(db, logger)
	accountServer := api.NewAccountServiceServer(accountService)

	walletService := services.NewWalletServiceImpl(db, logger, ks)
	walletServer := api.NewWalletServiceServer(walletService)

	authService := services.NewAuthServiceImpl(
		logger,
		accountService,
		opts.AppName,
		[]byte(opts.AuthConfig.SigningKey),
		int64(opts.AuthConfig.TokenExpirationTime),
	)
	authServer := api.NewAuthServiceServer(authService)

	discoveryService := services.NewDiscoveryServiceImpl(logger, walletService, ks, brokerContract)
	discoveryServiceServer := api.NewDiscoveryServiceServer(discoveryService)

	cryptoMessageService := services.NewCryptoMessageServiceImpl(logger, walletService)
	cryptoMessageServiceServer := api.NewCryptoMessageServiceServer(cryptoMessageService)

	userContractService := services.NewUserContractServiceImpl(logger, walletService, ks, userContract)
	userContractProxyServer := api.NewUserContractServiceServer(
		userContractService,
		userContract,
	)

	deviceContractService := services.NewDeviceContractServiceImpl(logger, walletService, ks, deviceContract)
	deviceContractProxyServer := api.NewDeviceContractServiceServer(
		deviceContractService,
		deviceContract,
	)

	productContractService := services.NewProductContractService(logger, walletService, ks, productContract)
	productContractProxyServer := api.NewProductContractServiceServer(
		productContractService,
		productContract,
	)

	brokerContractService := services.NewBrokerContractServiceImpl(logger, walletService, ks, brokerContract)
	brokerContractServer := api.NewBrokerContractServiceServer(
		brokerContractService,
		brokerContract,
	)

	negotiationContractService := services.NewNegotiationContractServiceImpl(
		logger,
		walletService,
		ks,
		negotiationContract,
	)
	negotiationContractServer := api.NewNegotiationContractServiceServer(
		negotiationContractService,
		negotiationContract,
	)

	biddingContractServer := api.NewBiddingContractServiceServer(logger, walletService, ks, ethClient)
	settlementContractServer := api.NewSettlementContractServiceServer(logger, walletService, ks, ethClient)

	tradingContractService := services.NewTradingContractServiceImpl(logger, walletService, ks, tradingContract)
	tradingContractServer := api.NewTradingContractServiceServer(
		tradingContractService,
		tradingContract,
	)

	grpcServer := initGrpcServer(authService, logger)
	api.RegisterAuthServiceServer(grpcServer, authServer)
	api.RegisterAccountServiceServer(grpcServer, accountServer)
	api.RegisterWalletServiceServer(grpcServer, walletServer)
	api.RegisterUserContractServiceServer(grpcServer, userContractProxyServer)
	api.RegisterDeviceContractServiceServer(grpcServer, deviceContractProxyServer)
	api.RegisterProductContractServiceServer(grpcServer, productContractProxyServer)
	api.RegisterBrokerContractServiceServer(grpcServer, brokerContractServer)
	api.RegisterNegotiationContractServiceServer(grpcServer, negotiationContractServer)
	api.RegisterBiddingContractServiceServer(grpcServer, biddingContractServer)
	api.RegisterTradingContractServiceServer(grpcServer, tradingContractServer)
	api.RegisterSettlementContractServiceServer(grpcServer, settlementContractServer)
	api.RegisterDiscoveryServiceServer(grpcServer, discoveryServiceServer)
	api.RegisterCryptoMessageServiceServer(grpcServer, cryptoMessageServiceServer)

	p := &proxy{
		opts:       opts,
		db:         db,
		logger:     logger,
		grpcServer: grpcServer,
		ethClient:  ethClient,
		running:    true,
		quit:       make(chan bool, 1),
	}

	return p, nil
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

func initDb(logger logrus.FieldLogger, opts options) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", opts.DatabaseConfig.Source)
	if err != nil {
		return nil, fmt.Errorf("open database %s: %w", opts.DatabaseConfig.Source, err)
	}
	db.SetLogger(logger)
	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(&model.Account{}, &model.Wallet{})
	return db, err
}

func initGrpcServer(authService services.AuthService, logger logrus.FieldLogger) *grpc.Server {
	entry := logrus.NewEntry(logger.(*logrus.Logger))
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(
				entry,
				grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
			),
			grpc_middleware_auth.StreamServerInterceptor(authService.AuthFunction()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(
				entry,
				grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
			),
			grpc_middleware_auth.UnaryServerInterceptor(authService.AuthFunction()),
		)),

	)
	grpc_logrus.ReplaceGrpcLogger(entry)
	return server
}

func (p *proxy) Run() error {
	addr := p.opts.Host + ":" + strconv.Itoa(p.opts.Port)

	p.logger.Infof("Starting proxy %s", version)

	if len(build) != 0 {
		p.logger.Infof("Built on %s", build)
	}

	if len(gitCommit) != 0 {
		p.logger.Infof("Git commit %s", gitCommit)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	p.logger.Infof("Proxy listening on %s", lis.Addr())

	p.receiveSignals()

	p.running = true
	return p.grpcServer.Serve(lis)
}

func (p *proxy) receiveSignals() {
	if p.opts.NoSig {
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go p.handleSignals(sigs)
}

func (p *proxy) handleSignals(sigs chan os.Signal) {
	select {
	case <-sigs:
		p.GraceFulStop()
		os.Exit(0)
	case <-p.quit:
		return
	}
}

func (p *proxy) GraceFulStop() {
	p.logger.Infof("Initiating graceful stop of proxy")
	if !p.running {
		return
	}
	close(p.quit)
	p.grpcServer.GracefulStop()
	p.ethClient.Close()
	err := p.db.Close()
	if err != nil {
		p.logger.Errorf("%+v", err)
	}
	p.running = false
	p.logger.Infof("Gracefully stopped proxy")
}
