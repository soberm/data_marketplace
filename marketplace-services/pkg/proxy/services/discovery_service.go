package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"marketplace-services/pkg/contracts"
	"math/big"
	"strings"
)

type DiscoveryService interface {
	SearchBroker(ctx context.Context, query *BrokerSearchQuery) (<-chan *contracts.Broker, <-chan error)
}

type BrokerSearchQuery struct {
	Name      string
	Locations []contracts.Location
}

type discoveryServiceImpl struct {
	logger         logrus.FieldLogger
	walletService  WalletService
	keyStore       *keystore.KeyStore
	brokerContract contracts.BrokerContract
}

func NewDiscoveryServiceImpl(
	logger logrus.FieldLogger,
	walletService WalletService,
	keyStore *keystore.KeyStore,
	brokerContract contracts.BrokerContract,
) *discoveryServiceImpl {
	return &discoveryServiceImpl{
		logger:         logger,
		walletService:  walletService,
		keyStore:       keyStore,
		brokerContract: brokerContract,
	}
}

func (s discoveryServiceImpl) SearchBroker(
	ctx context.Context,
	query *BrokerSearchQuery,
) (<-chan *contracts.Broker, <-chan error) {
	sink := make(chan *contracts.Broker)
	errc := make(chan error, 1)
	go func() {
		defer close(sink)
		defer close(errc)

		w, err := s.walletService.FindWalletByAuthenticatedAccount(ctx)
		if err != nil {
			errc <- fmt.Errorf("find wallet of authenticated proxy account: %w", err)
		}
		callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(w.Address)}

		count, err := s.brokerContract.CountBrokers(callOpts)
		if err != nil {
			errc <- fmt.Errorf("count brokers: %w", err)
			return
		}
		for i := uint64(0); i < count.Uint64(); i++ {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			default:
				broker, err := s.brokerContract.FindBrokerByIndex(callOpts, big.NewInt(int64(i)))
				if err != nil {
					errc <- fmt.Errorf("find broker by index %d: %w", i, err)
					return
				}
				if s.MatchSearchQuery(query, broker) {
					sink <- broker
				}
			}
		}
	}()
	return sink, errc
}

func (s discoveryServiceImpl) MatchSearchQuery(
	query *BrokerSearchQuery,
	broker *contracts.Broker,
) bool {
	matched := true
	if !strings.Contains(broker.Name, query.Name) {
		matched = false
	}
	if !contracts.ContainsLocation(query.Locations, broker.Location) {
		matched = false
	}
	return matched
}
