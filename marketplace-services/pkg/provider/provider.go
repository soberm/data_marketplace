package provider

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/api"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type provider struct {
	ctx                             context.Context
	opts                            options
	logger                          logrus.FieldLogger
	proxy                           *grpc.ClientConn
	authServiceClient               api.AuthServiceClient
	tradingContractServiceClient    api.TradingContractServiceClient
	deviceContractServiceClient     api.DeviceContractServiceClient
	settlementContractServiceClient api.SettlementContractServiceClient
	cryptoMessageServiceClient      api.CryptoMessageServiceClient
	brokerContractServiceClient     api.BrokerContractServiceClient
	simulators                      map[int]*sensorSimulator
}

func New(opt ...Option) (*provider, error) {
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

	proxy, err := grpc.Dial(opts.ProxyConfig.Address+":"+strconv.Itoa(opts.ProxyConfig.Port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial proxy: %w", err)
	}

	authServiceClient := api.NewAuthServiceClient(proxy)
	tradingServiceClient := api.NewTradingContractServiceClient(proxy)
	deviceServiceClient := api.NewDeviceContractServiceClient(proxy)
	settlementServiceClient := api.NewSettlementContractServiceClient(proxy)
	messageServiceClient := api.NewCryptoMessageServiceClient(proxy)
	brokerServiceClient := api.NewBrokerContractServiceClient(proxy)

	simulators := make(map[int]*sensorSimulator)
	for _, simulatorConfig := range opts.SimulationConfig.SimulatorConfigs {
		simulators[simulatorConfig.ID] = NewSensorSimulator(
			simulatorConfig.Min,
			simulatorConfig.Max,
			time.Duration(simulatorConfig.Frequency)*time.Second,
			time.Duration(simulatorConfig.Timeout)*time.Second,
		)
	}

	return &provider{
		opts:                            opts,
		logger:                          initLogger(opts),
		proxy:                           proxy,
		ctx:                             context.Background(),
		authServiceClient:               authServiceClient,
		tradingContractServiceClient:    tradingServiceClient,
		deviceContractServiceClient:     deviceServiceClient,
		settlementContractServiceClient: settlementServiceClient,
		cryptoMessageServiceClient:      messageServiceClient,
		brokerContractServiceClient:     brokerServiceClient,
		simulators:                      simulators,
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

func (p *provider) Run() error {
	p.logger.Infof("Starting provider %s", version)

	if len(build) != 0 {
		p.logger.Infof("Built on %s", build)
	}

	if len(gitCommit) != 0 {
		p.logger.Infof("Git commit %s", gitCommit)
	}

	for product, sim := range p.simulators {
		p.logger.Infof("Starting simulator for product %d", product)
		go sim.Simulate()
	}

	p.logger.Infof("Get access token for proxy %s", p.proxy.Target())

	getTokenRequest := &api.GetTokenRequest{
		Username: p.opts.ProxyConfig.Username,
		Password: []byte(p.opts.ProxyConfig.Password),
	}
	getTokenResponse, err := p.authServiceClient.GetToken(p.ctx, getTokenRequest)
	if err != nil {
		return err
	}

	md := metadata.Pairs("authorization", "bearer "+getTokenResponse.Token)
	ctx := metadata.NewOutgoingContext(p.ctx, md)

	p.receiveSignals()

	p.logger.Infof("Listen and serve trading requests")
	return p.listenAndServe(ctx)
}

func (p *provider) receiveSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go p.handleSignal(sigs)
}

func (p *provider) handleSignal(sigs chan os.Signal) {
	<-sigs
	err := p.proxy.Close()
	if err != nil {
		p.logger.Errorf("close proxy connection: %v", err)
	}
	p.logger.Infof("Shutdown provider")
	os.Exit(0)
}

func (p *provider) listenAndServe(ctx context.Context) error {
	products := p.providedProducts()
	stream, err := p.tradingContractServiceClient.WatchRequestedTradingEvent(
		ctx,
		&api.WatchRequestedTradingEventRequest{
			Products: products,
		},
	)
	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			return err
		}
		go func() {
			err := p.handleRequestedTradingEvent(ctx, response.Event)
			if err != nil {
				p.logger.Errorf("handle trading request: %+v", err)
			}
		}()
	}
}

func (p *provider) providedProducts() []uint64 {
	products := make([]uint64, 0, len(p.simulators))
	for product := range p.simulators {
		products = append(products, uint64(product))
	}
	return products
}

func (p *provider) handleRequestedTradingEvent(ctx context.Context, event *domain.RequestedTradingEvent) error {
	tradeId, err := p.acceptTradingRequest(ctx, event.Id)
	if err != nil {
		return err
	}

	findTradeResponse, err := p.tradingContractServiceClient.FindTradeById(ctx, &api.FindTradeByIdRequest{
		Id: tradeId,
	})
	if err != nil {
		return err
	}
	trade := findTradeResponse.Trade

	err = p.waitForDeposit(ctx, trade.SettlementContract)
	if err != nil {
		return err
	}

	findBrokerResponse, err := p.brokerContractServiceClient.FindBrokerByAddress(ctx, &api.FindBrokerByAddressRequest{
		Address: trade.Broker,
	})
	if err != nil {
		return err
	}
	broker := findBrokerResponse.Broker

	findDeviceResponse, err := p.deviceContractServiceClient.FindDeviceByAddress(ctx, &api.FindDeviceByAddressRequest{
		Address: trade.Consumer,
	})
	if err != nil {
		return err
	}
	device := findDeviceResponse.Device

	counter, err := p.pushMessages(ctx, trade, broker.HostAddr, device.PublicKey)
	if err != nil {
		return err
	}

	p.logger.Infof("Settle trade %d with counter %d", trade.Id, counter)
	_, err = p.settlementContractServiceClient.SettleTrade(ctx, &api.SettleTradeRequest{
		ContractAddress: trade.SettlementContract,
		Counter:         uint64(counter),
	})
	if err != nil {
		return err
	}
	p.logger.Infof("Finished trade %d", trade.Id)
	return nil
}

func (p *provider) acceptTradingRequest(ctx context.Context, id uint64) (uint64, error) {
	p.logger.Infof("Accept trading request %+v", id)
	stream, err := p.tradingContractServiceClient.WatchAcceptedTradingRequestEvent(
		ctx,
		&api.WatchAcceptedTradingRequestEventRequest{
			Ids: []uint64{id},
		},
	)
	if err != nil {
		return 0, err
	}

	_, err = p.tradingContractServiceClient.AcceptTradingRequest(ctx, &api.AcceptTradingRequestRequest{
		Id: id,
	})
	if err != nil {
		return 0, err
	}

	response, err := stream.Recv()
	if err != nil {
		return 0, err
	}

	return response.Event.Trade, nil
}

func (p *provider) waitForDeposit(ctx context.Context, contract string) error {
	p.logger.Infof("Waiting for deposit to contract %s", contract)
	stream, err := p.settlementContractServiceClient.WatchDepositedEvent(ctx, &api.WatchDepositedEventRequest{
		ContractAddress: contract,
	})
	if err != nil {
		return err
	}

	_, err = stream.Recv()
	if err != nil {
		return err
	}

	return err
}

func (p *provider) pushMessages(ctx context.Context, trade *domain.Trade, broker string, pubKey []byte) (int, error) {
	p.logger.Infof("Push messages for trade %d", trade.Id)
	sink := make(chan int)
	simulator, ok := p.simulators[int(trade.Product)]
	if !ok {
		return 0, fmt.Errorf("simulator for product %d not found", trade.Product)
	}

	startTime := time.Unix(int64(trade.StartTime), 0)
	waitTime := time.Until(startTime)
	p.logger.Infof("Waiting %s with transmission until trade %d starts", waitTime, trade.Id)
	<-time.After(waitTime)
	simulator.Attach(sink)

	endTime := time.Unix(int64(trade.EndTime), 0)
	stop := time.After(time.Until(endTime))

	counter := 0
	for {
		select {
		case m := <-sink:
			buf := new(bytes.Buffer)
			if err := binary.Write(buf, binary.LittleEndian, int64(m)); err != nil {
				p.logger.Errorf("measurement %d to binary: %w", m, err)
				break
			}
			_, err := p.cryptoMessageServiceClient.EncryptAndPushMessage(ctx, &api.EncryptAndPushMessageRequest{
				BrokerAddr: broker,
				PublicKey:  pubKey,
				Message: &domain.Message{
					TradeId: trade.Id,
					Payload: buf.Bytes(),
				},
			})
			if err != nil {
				p.logger.Errorf("push message for trade %d: %w", trade.Id, err)
				break
			}
			counter++
			p.logger.Infof("Pushed message with payload %d for trade %d", m, trade.Id)
		case <-stop:
			err := simulator.Detach(sink)
			if err != nil {
				p.logger.Errorf("detach sink: %w", trade.Id, err)
			}
			return counter, nil
		}
	}
}
