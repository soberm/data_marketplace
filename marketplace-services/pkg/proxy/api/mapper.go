package api

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/model"
	"marketplace-services/pkg/proxy/services"
	"math/big"
)

func AccountFromGrpcAccount(account *domain.Account) *model.Account {
	if account == nil {
		return nil
	}
	return &model.Account{
		Model: gorm.Model{
			ID: uint(account.Id),
		},
		Name:     account.Name,
		Password: account.Password,
		Role:     model.Role(account.Role),
	}
}

func AccountToGrpcAccount(account *model.Account) *domain.Account {
	if account == nil {
		return nil
	}
	return &domain.Account{
		Id:       uint64(account.ID),
		Name:     account.Name,
		Password: nil,
		Role:     domain.Role(account.Role),
		Wallet:   WalletToGrpcWallet(&account.Wallet),
	}
}

func WalletFromGrpcWallet(wallet *domain.Wallet) *model.Wallet {
	if wallet == nil {
		return nil
	}
	return &model.Wallet{
		Model: gorm.Model{
			ID: uint(wallet.Id),
		},
		AccountID:  uint(wallet.UserId),
		Passphrase: wallet.Passphrase,
		Address:    wallet.Address,
		FilePath:   wallet.FilePath,
		PublicKey:  wallet.PublicKey,
	}
}

func WalletToGrpcWallet(wallet *model.Wallet) *domain.Wallet {
	if wallet == nil {
		return nil
	}
	return &domain.Wallet{
		Id:         uint64(wallet.ID),
		UserId:     uint64(wallet.AccountID),
		Passphrase: "",
		Address:    wallet.Address,
		FilePath:   wallet.FilePath,
		PublicKey:  wallet.PublicKey,
	}
}

func TransactionToGrpcTransaction(tx *types.Transaction) *domain.Transaction {
	return &domain.Transaction{
		Hash:     tx.Hash().Hex(),
		Data:     tx.Data(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice().Uint64(),
		Value:    tx.Value().Int64(),
		Nonce:    tx.Nonce(),
	}
}

func UserFromGrpcUser(u *domain.User) *contracts.User {
	return &contracts.User{
		Addr:      common.HexToAddress(u.Address),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Company:   u.Company,
		Email:     u.Email,
	}
}

func UserToGrpcUser(u *contracts.User) *domain.User {
	return &domain.User{
		Address:   u.Addr.Hex(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Company:   u.Company,
		Email:     u.Email,
		Deleted:   u.Deleted,
		Brokers:   contracts.AddressesToHex(u.Brokers),
		Devices:   contracts.AddressesToHex(u.Devices),
	}
}

func BrokerFromGrpcBroker(broker *domain.Broker) *contracts.Broker {
	return &contracts.Broker{
		Addr:     common.HexToAddress(broker.Address),
		User:     common.HexToAddress(broker.User),
		Name:     broker.Name,
		HostAddr: broker.HostAddr,
		Location: contracts.Location(broker.Location),
	}
}

func BrokerToGrpcBroker(broker *contracts.Broker) *domain.Broker {
	return &domain.Broker{
		Address:  broker.Addr.Hex(),
		User:     broker.User.Hex(),
		Name:     broker.Name,
		HostAddr: broker.HostAddr,
		Location: domain.Location(broker.Location),
		Deleted:  broker.Deleted,
		Trades:   contracts.BigIntToInt64(broker.Trades),
	}
}

func DeviceFromGrpcDevice(device *domain.Device) *contracts.Device {
	return &contracts.Device{
		Addr:        common.HexToAddress(device.Address),
		User:        device.User,
		Name:        device.Name,
		Description: device.Description,
		PublicKey:   device.PublicKey[:],
		Deleted:     device.Deleted,
	}
}

func DeviceToGrpcDevice(device *contracts.Device) *domain.Device {
	return &domain.Device{
		Address:     device.Addr.Hex(),
		User:        device.User,
		Name:        device.Name,
		Description: device.Description,
		PublicKey:   device.PublicKey,
		Deleted:     device.Deleted,
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
		Cost:        product.Frequency.Int64(),
		Deleted:     product.Deleted,
	}
}

func BrokerSearchQueryFromGrpcBrokerSearchQuery(query *domain.BrokerSearchQuery) *services.BrokerSearchQuery {
	return &services.BrokerSearchQuery{
		Name:      query.Name,
		Locations: LocationsFromGrpcLocations(query.Locations),
	}
}

func LocationsFromGrpcLocations(locations []domain.Location) []contracts.Location {
	data := make([]contracts.Location, len(locations))
	for i := range data {
		data[i] = contracts.Location(locations[i])
	}
	return data
}

func BidFromGrpcBid(bid *domain.Bid) *contracts.Bid {
	return &contracts.Bid{
		Price:     big.NewInt(int64(bid.Price)),
		StartTime: big.NewInt(int64(bid.StartTime)),
		EndTime:   big.NewInt(int64(bid.EndTime)),
	}
}

func BidToGrpcBid(bid *contracts.Bid) *domain.Bid {
	return &domain.Bid{
		Price:     bid.Price.Uint64(),
		StartTime: bid.StartTime.Uint64(),
		EndTime:   bid.EndTime.Uint64(),
	}
}

func NegotiationRequestFromGrpcNegotiationRequest(negotiationRequest *domain.NegotiationRequest) *contracts.NegotiationRequest {
	return &contracts.NegotiationRequest{
		Id:       big.NewInt(int64(negotiationRequest.Id)),
		Product:  big.NewInt(int64(negotiationRequest.Product)),
		Consumer: common.HexToAddress(negotiationRequest.Consumer),
	}
}

func NegotiationRequestToGrpcNegotiationRequest(negotiationRequest *contracts.NegotiationRequest) *domain.NegotiationRequest {
	return &domain.NegotiationRequest{
		Id:       negotiationRequest.Id.Uint64(),
		Product:  negotiationRequest.Product.Uint64(),
		Consumer: negotiationRequest.Consumer.Hex(),
	}
}

func NegotiationFromGrpcNegotiation(negotiation *domain.Negotiation) *contracts.Negotiation {
	return &contracts.Negotiation{
		Id:              big.NewInt(int64(negotiation.Id)),
		Consumer:        common.HexToAddress(negotiation.Consumer),
		Product:         big.NewInt(int64(negotiation.Product)),
		BiddingContract: common.HexToAddress(negotiation.BiddingContract),
	}
}

func NegotiationToGrpcNegotiation(negotiation *contracts.Negotiation) *domain.Negotiation {
	return &domain.Negotiation{
		Id:              negotiation.Id.Uint64(),
		Consumer:        negotiation.Consumer.Hex(),
		Product:         negotiation.Product.Uint64(),
		BiddingContract: negotiation.BiddingContract.Hex(),
	}
}

func TradingRequestFromGrpcTradingRequest(tradingRequest *domain.TradingRequest) *contracts.TradingRequest {
	return &contracts.TradingRequest{
		Id:        big.NewInt(int64(tradingRequest.Id)),
		Product:   big.NewInt(int64(tradingRequest.Product)),
		Cost:      big.NewInt(int64(tradingRequest.Cost)),
		StartTime: big.NewInt(int64(tradingRequest.StartTime)),
		EndTime:   big.NewInt(int64(tradingRequest.EndTime)),
		Consumer:  common.HexToAddress(tradingRequest.Consumer),
		Broker:    common.HexToAddress(tradingRequest.Broker),
	}
}

func TradingRequestToGrpcTradingRequest(tradingRequest *contracts.TradingRequest) *domain.TradingRequest {
	return &domain.TradingRequest{
		Id:        tradingRequest.Id.Uint64(),
		Product:   tradingRequest.Product.Uint64(),
		Cost:      tradingRequest.Cost.Uint64(),
		StartTime: tradingRequest.StartTime.Uint64(),
		EndTime:   tradingRequest.EndTime.Uint64(),
		Consumer:  tradingRequest.Consumer.Hex(),
		Broker:    tradingRequest.Broker.Hex(),
	}
}

func TradeFromGrpcTrade(trade *domain.Trade) *contracts.Trade {
	return &contracts.Trade{
		Id:                 big.NewInt(int64(trade.Id)),
		Provider:           common.HexToAddress(trade.Provider),
		Consumer:           common.HexToAddress(trade.Consumer),
		Broker:             common.HexToAddress(trade.Broker),
		Product:            big.NewInt(int64(trade.Product)),
		StartTime:          big.NewInt(int64(trade.StartTime)),
		EndTime:            big.NewInt(int64(trade.EndTime)),
		Cost:               big.NewInt(int64(trade.Cost)),
		SettlementContract: common.HexToAddress(trade.SettlementContract),
	}
}

func TradeToGrpcTrade(trade *contracts.Trade) *domain.Trade {
	return &domain.Trade{
		Id:                 trade.Id.Uint64(),
		Provider:           trade.Provider.Hex(),
		Consumer:           trade.Consumer.Hex(),
		Broker:             trade.Broker.Hex(),
		Product:            trade.Product.Uint64(),
		StartTime:          trade.StartTime.Uint64(),
		EndTime:            trade.EndTime.Uint64(),
		Cost:               trade.Cost.Uint64(),
		SettlementContract: trade.SettlementContract.Hex(),
	}
}

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

func CounterToGrpcCounter(counter *contracts.Counter) *domain.Counter {
	return &domain.Counter{
		Value: counter.Value.Uint64(),
		Set:   counter.Set,
	}
}

func SettlementToGrpcSettlement(settlement *contracts.Settlement) *domain.Settlement {
	return &domain.Settlement{
		ActualCost: settlement.ActualCost.Uint64(),
		Provider:   settlement.Provider.Uint64(),
		Consumer:   settlement.Consumer.Uint64(),
		Broke:      settlement.Broker.Uint64(),
	}
}
