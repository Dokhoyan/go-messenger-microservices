package repository

//go:generate mockgen -destination=../repository/mocks/user_repository.go -package=mocks . UserRepository

import (
	"context"

	"github.com/Dokhoyan/common/pkg/filter"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
)

//postgresql
type UserRepository interface {
	Create(ctx context.Context, params *model.UserCreate) (int64, error)
	Get(ctx context.Context, filters filter.Filter) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error

}

type LogsRepository interface {
	Create(ctx context.Context, log model.Log) (int64, error)
	Get(ctx context.Context, id int64) (model.Log, error)
}