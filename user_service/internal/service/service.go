package service

//go:generate mockgen -destination=../service/mocks/user_service.go -package=mocks . UserService
//go:generate mockgen -destination=../service/mocks/auth_service.go -package=mocks . AuthService
//go:generate mockgen -destination=../service/mocks/access_service.go -package=mocks . AccessService

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, userParams *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, moodel *model.UserUpdate) (error)
	Delete(ctx context.Context, id int64) (error)
}

type AuthService interface {
	Login(context.Context, model.LoginDTO) (string, error)
	GetRefreshToken(context.Context, string) (string, error)
	GetAccessToken(context.Context, string) (string, error)
}

// AccessService - сервис проверки доступов
type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}