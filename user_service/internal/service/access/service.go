package access

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/utils"
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