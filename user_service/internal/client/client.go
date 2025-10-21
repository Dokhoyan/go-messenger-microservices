package client

import (
	"context"
)

// AuthService - сервис авторизации и аутентификации
type Auth interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

type KafkaProducer interface {
	Produce(ctx context.Context, topicName string, handler KafkaHandler) (error)
	Close() error
}

// Handler — функция, которая подготавливает данные для отправки в Kafka
type KafkaHandler interface {
	Data() (interface{}, error)
}