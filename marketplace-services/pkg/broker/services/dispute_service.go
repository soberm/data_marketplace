package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"marketplace-services/pkg/contracts"
	"marketplace-services/pkg/contracts/bindings"
	"math/big"
)

type DisputeService interface {
	ResolveDisputes(ctx context.Context) error
}

type disputeServiceImpl struct {
	logger          logrus.FieldLogger
	keyStore        *keystore.KeyStore
	ethClient       *ethclient.Client
	tradingContract contracts.TradingContract
	messageService  MessageService
	account         string
	passphrase      string
}

func NewDisputeServiceImpl(
	logger logrus.FieldLogger,
	keystore *keystore.KeyStore,
	ethClient *ethclient.Client,
	tradingContract contracts.TradingContract,
	messageService MessageService,
	account string,
	passphrase string,
) *disputeServiceImpl {
	return &disputeServiceImpl{
		logger:          logger,
		keyStore:        keystore,
		ethClient:       ethClient,
		tradingContract: tradingContract,
		messageService:  messageService,
		account:         account,
		passphrase:      passphrase,
	}
}

func (d *disputeServiceImpl) ResolveDisputes(ctx context.Context) error {
	sink := make(chan *bindings.TradingContractCreatedTrade)
	defer close(sink)

	sub, err := d.tradingContract.WatchCreatedTrade(
		&bind.WatchOpts{
			Context: ctx,
		},
		sink,
		contracts.HexToAddresses([]string{d.account}),
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			callOpts := &bind.CallOpts{Context: ctx, From: common.HexToAddress(d.account)}
			trade, err := d.tradingContract.FindTradeById(callOpts, event.TradeId)
			if err != nil {
				return fmt.Errorf("find trade by id %d: %w", event.TradeId, err)
			}
			go func() {
				err := d.resolveDispute(ctx, event.TradeId.Uint64(), trade.SettlementContract)
				if err != nil {
					d.logger.Errorf("resolve dispute with contract %s: %v", trade.SettlementContract.Hex(), err)
				}
			}()
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (d *disputeServiceImpl) resolveDispute(ctx context.Context, tradeId uint64, address common.Address) error {
	settlementContract, err := contracts.NewSettlementContractImpl(address, d.ethClient)
	if err != nil {
		return fmt.Errorf("new settlement contract with address %s: %w", address.Hex(), err)
	}
	sink := make(chan *bindings.SettlementContractDispute)
	defer close(sink)

	sub, err := settlementContract.WatchDisputeEvent(
		&bind.WatchOpts{
			Context: ctx,
		},
		sink,
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	select {
	case event := <-sink:
		d.logger.Infof(
			"Handle dispute of contract %s with provider count %d and consumer count %d",
			address.Hex(),
			event.ProviderCounter,
			event.ConsumerCounter,
		)
		account := accounts.Account{Address: common.HexToAddress(d.account)}
		transactOpts, err := bind.NewKeyStoreTransactor(
			d.keyStore,
			account,
		)
		if err != nil {
			return fmt.Errorf("new keystore transactor: %w", err)
		}
		counter := d.messageService.FindCounter(tradeId)

		if err := d.keyStore.Unlock(account, d.passphrase); err != nil {
			return fmt.Errorf("unlock account %s: %w", account.Address.Hex(), err)
		}
		defer func() {
			if err := d.keyStore.Lock(account.Address); err != nil {
				d.logger.Warnf("lock account %s: %w", account.Address.Hex(), err)
			}
		}()
		_, err = settlementContract.ResolveDispute(transactOpts, big.NewInt(int64(counter)))
		if err != nil {
			return fmt.Errorf("settle trade with contract %s and counter %d: %w", address.Hex(), 0, err)
		}
	case err = <-sub.Err():
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
