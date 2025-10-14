package config

import (
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
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

type UserConfig interface {
	Address() string
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
