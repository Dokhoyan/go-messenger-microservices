package auth

import (
	"github.com/Dokhoyan/common/pkg/storage"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
)


type serv struct {
	redis     storage.Redis
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

// NewService - создает экземпляр сервиса авторизации
func NewService(redis storage.Redis, userRepo  repository.UserRepository, jwtConfig config.JWTConfig) service.AuthService {
	return &serv{
		redis:      redis,
		userRepo: userRepo,
		jwtConfig:  jwtConfig,
	}
}