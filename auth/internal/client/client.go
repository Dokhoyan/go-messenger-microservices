package client

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client/kafka/consumer"
)

type KafkaConsumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}