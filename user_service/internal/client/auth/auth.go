package auth

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client"
	accesspb "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/auth/proto"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type auth struct {
	client     accesspb.AccessV1Client
	authConfig config.AuthConfig
	conn       *grpc.ClientConn 
}

// NewAuthClient - создает новый инстанс подключения к сервису auth
func NewAuthClient(authConfig config.AuthConfig) (client.Auth, error) {
	conn, err := grpc.Dial(authConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	cl := accesspb.NewAccessV1Client(conn)
	return &auth{
		client:     cl,
		conn:       conn, 
		authConfig: authConfig,
	}, nil
}

func (a *auth) Check(ctx context.Context, endpoint string) (bool, error) {
	if _, err := a.client.Check(ctx, &accesspb.CheckRequest{
		EndpointAddress: endpoint,
	}); err != nil {
		return false, fmt.Errorf("accessClient.Check: %w", err)
	}

	return true, nil
}

func (a *auth) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}