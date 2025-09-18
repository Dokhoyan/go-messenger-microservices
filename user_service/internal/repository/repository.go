package repository

//go:generate mockgen -destination=../repository/mocks/user_repository.go -package=mocks . UserRepository

import (
	"context"
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/go-redis/redis"
)

//postgresql
type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}

//redis
type AuthRepository interface {
	Ping() error
	Close() error
	Get(string) *redis.StringCmd
	Set(string, interface{}, time.Duration) *redis.StatusCmd
}