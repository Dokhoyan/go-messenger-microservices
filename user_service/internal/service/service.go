package service

//go:generate mockgen -destination=../service/mocks/user_service.go -package=mocks . UserService

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