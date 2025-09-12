package repository

//go:generate mockgen -destination=../repository/mocks/user_repository.go -package=mocks . UserRepository

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}