package auth

import (
	"github.com/Dokhoyan/common/pkg/storage"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"

	//"github.com/Dokhoyan/go-messenger-microservices/auth/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
)


type serv struct {
	redis     storage.Redis
	userClient client.UserService
	jwtConfig config.JWTConfig
}

// NewService - создает экземпляр сервиса авторизации
func NewService(redis storage.Redis, userClient client.UserService, jwtConfig config.JWTConfig) service.AuthService {
	return &serv{
		redis:      redis,
		userClient: userClient,
		jwtConfig:  jwtConfig,
	}
}