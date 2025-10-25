package config

import (
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type JWTConfig interface {
	RefreshSecretKey() []byte
	RefreshExpirationTime() time.Duration
	AccessSecretKey() []byte
	AccessExpirationTime() time.Duration
}

type RedisConfig interface {
	Address() string
	Password() string
	RoutesAccesses() map[string][]model.UserRole
}


// GRPCConfig - конфиг gRPC
type GRPCConfig interface {
	Address() string
}

type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

type PGConfig interface {
	DSN() string
}