package services

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Message struct {
	TradeId uint64
	Payload []byte
}

type MessageService interface {
	PushMessage(ctx context.Context, message *Message) error
	PullMessage(ctx context.Context, tradeId uint64) (*Message, error)
	FindCounter(tradeId uint64) uint64
}

type messageServiceImpl struct {
	logger   logrus.FieldLogger
	queues   map[uint64]chan *Message
	counters map[uint64]uint64
}

func NewMessageServiceImpl(
	logger logrus.FieldLogger,
) *messageServiceImpl {
	return &messageServiceImpl{
		logger:   logger,
		queues:   make(map[uint64]chan *Message),
		counters: make(map[uint64]uint64),
	}
}

func (s messageServiceImpl) PushMessage(ctx context.Context, message *Message) error {
	queue, ok := s.queues[message.TradeId]
	if !ok {
		queue = make(chan *Message, 100)
		s.queues[message.TradeId] = queue
	}

	select {
	case queue <- message:
		s.logger.Infof("Enqueued message for trade %d", message.TradeId)
		s.counters[message.TradeId] = s.counters[message.TradeId] + 1
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s messageServiceImpl) PullMessage(ctx context.Context, tradeId uint64) (*Message, error) {
	queue, ok := s.queues[tradeId]
	if !ok {
		queue = make(chan *Message, 100)
		s.queues[tradeId] = queue
	}

	select {
	case message := <-queue:
		s.logger.Infof("Dequeued message for trade %d", message.TradeId)
		return message, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s messageServiceImpl) FindCounter(tradeId uint64) uint64 {
	return s.counters[tradeId]
}
