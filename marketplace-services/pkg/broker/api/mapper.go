package api

import (
	"github.com/ethereum/go-ethereum/common"
	"marketplace-services/pkg/broker/services"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/domain"
	"math/big"
)

func MessageFromGrpcMessage(message *domain.Message) *services.Message {
	return &services.Message{
		TradeId: message.TradeId,
		Payload: message.Payload,
	}
}

func MessageToGrpcMessage(message *services.Message) *domain.Message {
	return &domain.Message{
		TradeId: message.TradeId,
		Payload: message.Payload,
	}
}

func ProductSearchQueryFromGrpcProductSearchQuery(query *domain.ProductSearchQuery) *services.ProductSearchQuery {
	return &services.ProductSearchQuery{
		DataType:     query.DataType,
		MinCost:      query.MinCost,
		MaxCost:      query.MaxCost,
		MinFrequency: query.MinFrequency,
		MaxFrequency: query.MaxFrequency,
	}
}

func ProductSearchQueryToGrpcProductSearchQuery(query *services.ProductSearchQuery) *domain.ProductSearchQuery {
	return &domain.ProductSearchQuery{
		DataType:     query.DataType,
		MinCost:      query.MinCost,
		MaxCost:      query.MaxCost,
		MinFrequency: query.MinFrequency,
		MaxFrequency: query.MaxFrequency,
	}
}

func ProductFromGrpcProduct(product *domain.Product) *contracts.Product {
	return &contracts.Product{
		Id:          big.NewInt(product.Id),
		Device:      common.HexToAddress(product.Device),
		Name:        product.Name,
		Description: product.Description,
		DataType:    product.DataType,
		Frequency:   big.NewInt(product.Frequency),
		Cost:        big.NewInt(product.Cost),
		Deleted:     product.Deleted,
	}
}

func ProductToGrpcProduct(product *contracts.Product) *domain.Product {
	return &domain.Product{
		Id:          product.Id.Int64(),
		Device:      product.Device.Hex(),
		Name:        product.Name,
		Description: product.Description,
		DataType:    product.DataType,
		Frequency:   product.Frequency.Int64(),
		Cost:        product.Cost.Int64(),
		Deleted:     product.Deleted,
	}
}
