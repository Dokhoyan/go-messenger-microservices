package repository

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, userinfo *model.UserAuthData) (int64, error)
	Get(ctx context.Context, userName string) (*model.User, error)
}