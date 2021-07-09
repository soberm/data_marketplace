package consumer

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/api"
	"os"
	"strconv"
	"time"
)

type consumer struct {
	opts                            options
	ctx                             context.Context
	logger                          logrus.FieldLogger
	proxy                           *grpc.ClientConn
	authServiceClient               api.AuthServiceClient
	walletServiceClient             api.WalletServiceClient
	discoveryServiceClient          api.DiscoveryServiceClient
	tradingContractServiceClient    api.TradingContractServiceClient
	settlementContractServiceClient api.SettlementContractServiceClient
	cryptoMessageServiceClient      api.CryptoMessageServiceClient
}

func New(opt ...Option) (*consumer, error) {
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

	proxy, err := grpc.Dial(opts.ProxyConfig.Address+":"+strconv.Itoa(opts.ProxyConfig.Port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial proxy: %w", err)
	}
	authService := api.NewAuthServiceClient(proxy)
	walletService := api.NewWalletServiceClient(proxy)
	discoveryService := api.NewDiscoveryServiceClient(proxy)
	tradingService := api.NewTradingContractServiceClient(proxy)
	settlementService := api.NewSettlementContractServiceClient(proxy)
	messageService := api.NewCryptoMessageServiceClient(proxy)
	return &consumer{
		opts:                            opts,
		ctx:                             context.Background(),
		logger:                          initLogger(opts),
		proxy:                           proxy,
		authServiceClient:               authService,
		walletServiceClient:             walletService,
		discoveryServiceClient:          discoveryService,
		tradingContractServiceClient:    tradingService,
		settlementContractServiceClient: settlementService,
		cryptoMessageServiceClient:      messageService,
	}, nil
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

func (c *consumer) Run() error {

	c.logger.Infof("Starting consumer %s", version)

	if len(build) != 0 {
		c.logger.Infof("Built on %s", build)
	}

	if len(gitCommit) != 0 {
		c.logger.Infof("Git commit %s", gitCommit)
	}

	proxyConfig := c.opts.ProxyConfig
	searchConfig := c.opts.SearchConfig

	c.logger.Infof("Get access token for proxy %s", c.proxy.Target())

	getTokenRequest := &api.GetTokenRequest{
		Username: proxyConfig.Username,
		Password: []byte(proxyConfig.Password),
	}
	getTokenResponse, err := c.authServiceClient.GetToken(c.ctx, getTokenRequest)
	if err != nil {
		return err
	}

	md := metadata.Pairs("authorization", "bearer "+getTokenResponse.Token)
	ctx := metadata.NewOutgoingContext(c.ctx, md)

	c.logger.Infof("Searching for a broker in location %s", domain.Location(searchConfig.BrokerLocation))

	brokerStream, err := c.discoveryServiceClient.SearchBroker(ctx, &api.SearchBrokerRequest{
		Query: &domain.BrokerSearchQuery{
			Locations: []domain.Location{domain.Location(searchConfig.BrokerLocation)},
		},
	})
	if err != nil {
		return err
	}
	searchBrokerResponse, err := brokerStream.Recv()
	if err != nil {
		return err
	}

	c.logger.Infof(
		"Search for a product with type %s in price range [%d-%d] and frequency range [%d-%d]",
		searchConfig.DataType,
		searchConfig.MinCost,
		searchConfig.MaxCost,
		searchConfig.MinFrequency,
		searchConfig.MaxFrequency,
	)

	productStream, err := c.discoveryServiceClient.SearchProductWithBroker(ctx, &api.SearchProductWithBrokerRequest{
		BrokerAddr: searchBrokerResponse.Broker.HostAddr,
		Query: &domain.ProductSearchQuery{
			DataType:     searchConfig.DataType,
			MinCost:      searchConfig.MinCost,
			MaxCost:      searchConfig.MaxCost,
			MinFrequency: searchConfig.MinFrequency,
			MaxFrequency: searchConfig.MaxFrequency,
		},
	})
	if err != nil {
		return err
	}

	searchProductResponse, err := productStream.Recv()
	if err != nil {
		return err
	}

	c.logger.Infof("Request to buy product %d", searchProductResponse.Product.Id)

	watchRequestedTradingStream, err := c.tradingContractServiceClient.WatchRequestedTradingEvent(
		ctx,
		&api.WatchRequestedTradingEventRequest{
			Requesters: []string{proxyConfig.Account},
			Products:   []uint64{uint64(searchProductResponse.Product.Id)},
		})

	now := time.Now()
	startTime := now.Add(30 * time.Second)
	endTime := startTime.Add(60 * time.Second)
	_, err = c.tradingContractServiceClient.RequestTrading(ctx, &api.RequestTradingRequest{
		Product:   uint64(searchProductResponse.Product.Id),
		Broker:    searchBrokerResponse.Broker.Address,
		StartTime: uint64(startTime.Unix()),
		EndTime:   uint64(endTime.Unix()),
	})
	if err != nil {
		return err
	}

	watchRequestTradingResponse, err := watchRequestedTradingStream.Recv()
	if err != nil {
		return err
	}

	c.logger.Infof("Wait until provider accepts trading request %d", watchRequestTradingResponse.Event.Id)

	watchAcceptedTradingStream, err := c.tradingContractServiceClient.WatchAcceptedTradingRequestEvent(
		ctx,
		&api.WatchAcceptedTradingRequestEventRequest{
			Ids: []uint64{watchRequestTradingResponse.Event.Id},
		})
	if err != nil {
		return err
	}

	watchAcceptedTradingResponse, err := watchAcceptedTradingStream.Recv()
	if err != nil {
		return err
	}

	c.logger.Infof("Find trade for the accepted trading request %d", watchRequestTradingResponse.Event.Id)

	findTradeResponse, err := c.tradingContractServiceClient.FindTradeById(ctx, &api.FindTradeByIdRequest{
		Id: watchAcceptedTradingResponse.Event.Trade,
	})
	if err != nil {
		return err
	}

	c.logger.Infof("Deposit the funds for trade %d", findTradeResponse.Trade.Id)

	_, err = c.settlementContractServiceClient.Deposit(ctx, &api.DepositRequest{
		ContractAddress: findTradeResponse.Trade.SettlementContract,
		Value:           100000,
	})
	if err != nil {
		return err
	}

	waitTime := time.Until(startTime)
	c.logger.Infof("Waiting %s with transmission until trade %d starts", waitTime, findTradeResponse.Trade.Id)

	stop := time.After(time.Until(endTime))
	deadlineContext, _ := context.WithDeadline(ctx, endTime)

	counter := 0
	for {
		select {
		case <-stop:
			c.logger.Infof("Settle trade with counter %d", counter)
			_, err = c.settlementContractServiceClient.SettleTrade(ctx, &api.SettleTradeRequest{
				ContractAddress: findTradeResponse.Trade.SettlementContract,
				Counter:         uint64(counter),
			})
			c.logger.Infof("Finished trade %d", findTradeResponse.Trade.Id)
			c.logger.Infof("Shutdown consumer")
			return err
		default:
			pullMessageResponse, err := c.cryptoMessageServiceClient.DecryptAndPullMessage(deadlineContext, &api.DecryptAndPullMessageRequest{
				BrokerAddr: searchBrokerResponse.Broker.HostAddr,
				TradeId:    findTradeResponse.Trade.Id,
			})
			if err != nil {
				break
			}

			counter++
			payload := pullMessageResponse.Message.Payload
			value := int64(binary.LittleEndian.Uint64(payload))
			if err != nil {
				break
			}
			c.logger.Printf("Received value %d", value)
		}
	}
}
