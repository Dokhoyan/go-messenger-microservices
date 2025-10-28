package config

import (
	"errors"
	"os"
)

type GatewayConfig struct {
	Host string
	Port string
}

type ServicesConfig struct {
	AuthServiceURL     string
	ChatServiceURL     string
	UserServiceURL     string
	UserServiceHTTPURL string
}

func Load() (*GatewayConfig, *ServicesConfig, error) {
	cfg := &GatewayConfig{
		Host: getEnv("GATEWAY_HOST", "0.0.0.0"),
		Port: getEnv("GATEWAY_PORT", "8080"),
	}

	if cfg.Port == "" {
		return nil, nil, errors.New("gateway port not found")
	}

	servicesConfig := &ServicesConfig{
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://auth-service:50051"),
		ChatServiceURL:     getEnv("CHAT_SERVICE_URL", "http://chat-service:50052"),
		UserServiceURL:     getEnv("USER_SERVICE_URL", "http://user-service:50053"),
		UserServiceHTTPURL: getEnv("USER_SERVICE_HTTP_URL", "http://user-service:8090"),
	}

	return cfg, servicesConfig, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
