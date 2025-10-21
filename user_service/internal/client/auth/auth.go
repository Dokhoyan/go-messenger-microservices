package auth

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client"
	accessv1 "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/auth/proto"
)

type authclient struct {
	accessClient accessv1.AccessV1Client
}

func NewClient(cl accessv1.AccessV1Client) client.Auth {
	return &authclient{
		accessClient: cl,
	}
}

func (a *authclient) Check(ctx context.Context, endpoint string) (bool, error) {
	if _, err := a.accessClient.Check(ctx, &accessv1.CheckRequest{
		EndpointAddress: endpoint,
	}); err != nil {
		return false, fmt.Errorf("accessClient.Check: %w", err)
	}

	return true, nil
}