package access

import (
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/utils"
)

type serv struct {
	jwtConfig     config.JWTConfig
	accessChecker utils.AccessChecker
}

// NewService - создает новый экземпляр сервиса проверки доступа
func NewService(jwtConfig config.JWTConfig, checker utils.AccessChecker) service.AccessService {
	return &serv{
		jwtConfig:     jwtConfig,
		accessChecker: checker,
	}
}