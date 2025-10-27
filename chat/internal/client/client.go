package client

import "context"

type AuthService interface {
	Check(ctx context.Context, endpoint string) error
}