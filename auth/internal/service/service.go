package service

//go:generate mockgen -destination=../service/mocks/auth_service.go -package=mocks . AuthService
//go:generate mockgen -destination=../service/mocks/access_service.go -package=mocks . AccessService

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
)

type AuthService interface {
	Login(context.Context, model.LoginDTO) (string, error)
	GetRefreshToken(context.Context, string) (string, error)
	GetAccessToken(context.Context, string) (string, error)
}

// AccessService - сервис проверки доступов
type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}