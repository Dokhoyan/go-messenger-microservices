package service

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64)(*model.User, error)
}