package config

import (
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
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

type RedisConfig interface {
	Address() string
	Password() string
	RoutesAccesses() map[string][]model.UserRole
}

// SwaggerConfig - конфиг swagger
type SwaggerConfig interface {
	Address() string
}

// HTTPConfig - конфиг http
type HTTPConfig interface {
	Address() string
}

// GRPCConfig - конфиг gRPC
type GRPCConfig interface {
	Address() string
}

// PGConfig - конфиг Postgres
type PGConfig interface {
	DSN() string
}

type JWTConfig interface {
	RefreshSecretKey() []byte
	RefreshExpirationTime() time.Duration
	AccessSecretKey() []byte
	AccessExpirationTime() time.Duration
}

type PrometheusConfig interface {
	Address() string
}

type AuthConfig interface {
	Address() string
}

type KafkaProducerConfig interface {
	Brokers() []string
	Config() *sarama.Config
}