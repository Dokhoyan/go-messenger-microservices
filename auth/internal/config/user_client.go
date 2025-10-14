package config

import (
	"errors"
	"net"
	"os"
)

const (
	userClientHost = "USER_HOST"
	userClientPort = "USER_PORT"
)

type userConfig struct {
	host string
	port string
}

// NewAuthConfig - создает конфиг gRPC
func NewUserConfig() (UserConfig, error) {
	host := os.Getenv(userClientHost)
	if len(host) == 0 {
		return nil, errors.New("auth host not found in environments")
	}

	port := os.Getenv(userClientPort)
	if len(host) == 0 {
		return nil, errors.New("auth port not found in environments")
	}

	return userConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg userConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}