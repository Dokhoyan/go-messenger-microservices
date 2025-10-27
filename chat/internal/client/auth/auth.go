package auth

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/client"
	accesspb "github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/client/auth/proto"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type auth struct {
	client     accesspb.AccessV1Client
	authConfig config.AuthConfig
}

// NewAuthClient - создает новый инстанс подключения к сервису auth
func NewAuthClient(authConfig config.AuthConfig) (client.AuthService, error) {
	conn, err := grpc.Dial(authConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	cl := accesspb.NewAccessV1Client(conn)
	return &auth{
		client:     cl,
		authConfig: authConfig,
	}, nil
}

func (a *auth) Check(ctx context.Context, endpoint string) error {
	_, err := a.client.Check(ctx, &accesspb.CheckRequest{
		EndpointAddress: endpoint,
	})

	return err
}