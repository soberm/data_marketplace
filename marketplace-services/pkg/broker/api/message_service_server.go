package api

import (
	"context"
	"marketplace-services/pkg/broker/services"
)

type messageServiceServer struct {
	UnimplementedMessageServiceServer
	messageService services.MessageService
}

func NewMessageServiceServer(messageService services.MessageService) *messageServiceServer {
	return &messageServiceServer{messageService: messageService}
}

func (s *messageServiceServer) PushMessage(ctx context.Context, req *PushMessageRequest) (*PushMessageResponse, error) {
	err := s.messageService.PushMessage(ctx, MessageFromGrpcMessage(req.Message))
	if err != nil {
		return &PushMessageResponse{}, err
	}
	return &PushMessageResponse{}, err
}

func (s *messageServiceServer) PullMessage(ctx context.Context, req *PullMessageRequest) (*PullMessageResponse, error) {
	msg, err := s.messageService.PullMessage(ctx, req.TradeId)
	if err != nil {
		return &PullMessageResponse{}, err
	}
	return &PullMessageResponse{Message: MessageToGrpcMessage(msg)}, err
}
