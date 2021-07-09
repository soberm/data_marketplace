package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"marketplace-services/pkg/broker/api"
	"marketplace-services/pkg/domain"
	"math/rand"
	"time"
)

type Message struct {
	TradeId uint64
	Payload []byte
}

type CryptoMessageService interface {
	EncryptAndPushMessage(ctx context.Context, brokerAddress string, publicKey []byte, msg *Message) error
	DecryptAndPullMessage(ctx context.Context, brokerAddress string, tradeId uint64) (*Message, error)
}

type cryptoMessageServiceImpl struct {
	logger        logrus.FieldLogger
	walletService WalletService
}

func NewCryptoMessageServiceImpl(logger logrus.FieldLogger, walletService WalletService, ) *cryptoMessageServiceImpl {
	return &cryptoMessageServiceImpl{logger: logger, walletService: walletService}
}

func (c *cryptoMessageServiceImpl) EncryptAndPushMessage(ctx context.Context, brokerAddress string, publicKey []byte, msg *Message) error {

	conn, err := grpc.Dial(brokerAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("dial grpc %s: %w", brokerAddress, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			c.logger.Errorf("%+v", err)
		}
	}()

	pubKey, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		return fmt.Errorf("unmarshal pubkey: %w", err)
	}
	c.logger.Infof("%+v", msg)
	encryptedPayload, err := ecies.Encrypt(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		ecies.ImportECDSAPublic(pubKey),
		msg.Payload,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("encrypt payload: %w", err)
	}

	messageService := api.NewMessageServiceClient(conn)
	_, err = messageService.PushMessage(ctx, &api.PushMessageRequest{
		Message: &domain.Message{TradeId: msg.TradeId, Payload: encryptedPayload},
	})

	return err
}

func (c *cryptoMessageServiceImpl) DecryptAndPullMessage(ctx context.Context, brokerAddress string, tradeId uint64) (*Message, error) {
	conn, err := grpc.Dial(brokerAddress, grpc.WithInsecure())
	if err != nil {
		return &Message{}, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			c.logger.Errorf("%+v", err)
		}
	}()

	messageService := api.NewMessageServiceClient(conn)
	response, err := messageService.PullMessage(ctx, &api.PullMessageRequest{
		TradeId: tradeId,
	})
	if err != nil {
		return &Message{}, err
	}
	msg := response.Message

	key, err := c.walletService.FindKeyByAuthenticatedAccount(ctx)
	if err != nil {
		return &Message{}, fmt.Errorf("find key of authenticated proxy account: %w", err)
	}
	privateKey := ecies.ImportECDSA(key.PrivateKey)

	decryptedPayload, err := privateKey.Decrypt(msg.Payload, nil, nil)
	if err != nil {
		return &Message{}, fmt.Errorf("decrypt payload: %w", err)
	}

	return &Message{TradeId: msg.TradeId, Payload: decryptedPayload}, nil
}
