package api

import (
	"context"
	"marketplace-services/pkg/proxy/services"
)

type cryptoMessageServiceServer struct {
	UnimplementedCryptoMessageServiceServer
	cryptoMessageService services.CryptoMessageService
}

func NewCryptoMessageServiceServer(cryptoMessageService services.CryptoMessageService) *cryptoMessageServiceServer {
	return &cryptoMessageServiceServer{cryptoMessageService: cryptoMessageService}
}

func (s *cryptoMessageServiceServer) EncryptAndPushMessage(
	ctx context.Context,
	req *EncryptAndPushMessageRequest,
) (*EncryptAndPushMessageResponse, error) {
	return &EncryptAndPushMessageResponse{}, s.cryptoMessageService.EncryptAndPushMessage(
		ctx,
		req.BrokerAddr,
		req.PublicKey,
		MessageFromGrpcMessage(req.Message),
	)
}

func (s *cryptoMessageServiceServer) DecryptAndPullMessage(
	ctx context.Context,
	req *DecryptAndPullMessageRequest,
) (*DecryptAndPullMessageResponse, error) {
	msg, err := s.cryptoMessageService.DecryptAndPullMessage(ctx, req.BrokerAddr, req.TradeId)
	return &DecryptAndPullMessageResponse{Message: MessageToGrpcMessage(msg)}, err
}
