package authservice

import (
	"github.com/Dokhoyan/common/pkg/storage"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
)


type serv struct {
	redis     storage.Redis
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

// NewService - создает экземпляр сервиса авторизации
func NewService(redis storage.Redis, repo repository.UserRepository, jwtConfig config.JWTConfig) service.AuthService {
	return &serv{
		redis:     redis,
		userRepo:  repo,
		jwtConfig: jwtConfig,
	}
}