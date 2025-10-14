package client

import (
	"context"

	user_v1 "github.com/Dokhoyan/go-messenger-microservices/auth/internal/client/user/proto"

)

type UserService interface {
	GetUserAuthData(ctx context.Context, username string) (*user_v1.UserAuthData, error)
	Close() error
}