package auth

import (
	"context"
	"fmt"

	access_v1 "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/grpc/auth/proto"
)

var _ Client = (*client)(nil)

type Client interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

type client struct {
	accessClient access_v1.AccessV1Client
}

func NewClient(cl access_v1.AccessV1Client) *client {
	return &client{
		accessClient: cl,
	}
}

func (c *client) Check(ctx context.Context, endpoint string) (bool, error) {
	if _, err := c.accessClient.Check(ctx, &access_v1.CheckRequest{
		EndpointAddress: endpoint,
	}); err != nil {
		return false, fmt.Errorf("accessClient.Check: %w", err)
	}

	return true, nil
}