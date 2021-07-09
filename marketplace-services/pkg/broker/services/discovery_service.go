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
)

type ProductSearchQuery struct {
	DataType     string
	MinCost      uint64
	MaxCost      uint64
	MinFrequency uint64
	MaxFrequency uint64
}

type DiscoveryService interface {
	SearchProduct(ctx context.Context, query *ProductSearchQuery) (<-chan *contracts.Product, <-chan error)
}

type discoveryServiceImpl struct {
	logger          logrus.FieldLogger
	productContract contracts.ProductContract
	keyStore        *keystore.KeyStore
}

func NewDiscoveryServiceImpl(
	logger logrus.FieldLogger,
	productContract contracts.ProductContract,
	keyStore *keystore.KeyStore,
) *discoveryServiceImpl {
	return &discoveryServiceImpl{logger: logger, productContract: productContract, keyStore: keyStore}
}

func (s discoveryServiceImpl) SearchProduct(
	ctx context.Context,
	query *ProductSearchQuery,
) (<-chan *contracts.Product, <-chan error) {
	sink := make(chan *contracts.Product)
	errc := make(chan error, 1)
	callOpts := &bind.CallOpts{Context: ctx, From: common.BytesToAddress(nil)}
	go func() {
		defer close(sink)
		defer close(errc)

		count, err := s.productContract.CountProducts(callOpts)
		if err != nil {
			errc <- fmt.Errorf("count products: %w", err)
			return
		}
		for i := uint64(0); i < count.Uint64(); i++ {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			default:
				product, err := s.productContract.FindProductByIndex(callOpts, big.NewInt(int64(i)))
				if err != nil {
					errc <- fmt.Errorf("find product by index %d: %w", i, err)
					return
				}
				if s.MatchSearchQuery(query, product) {
					sink <- product
				}
			}
		}
	}()
	return sink, errc
}

func (s discoveryServiceImpl) MatchSearchQuery(
	query *ProductSearchQuery,
	product *contracts.Product,
) bool {
	matched := true
	if product.DataType != query.DataType && query.DataType != "" {
		matched = false
	}
	if product.Cost.Uint64() < query.MinCost {
		matched = false
	}
	if product.Cost.Uint64() > query.MaxCost && query.MaxCost != 0 {
		matched = false
	}
	if product.Frequency.Uint64() < query.MinFrequency {
		matched = false
	}
	if product.Frequency.Uint64() > query.MaxFrequency && query.MaxFrequency != 0 {
		matched = false
	}

	return matched
}
