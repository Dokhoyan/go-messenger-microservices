package user

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client"
	user_v1 "github.com/Dokhoyan/go-messenger-microservices/auth/internal/client/user/proto"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type user struct {
	userConfig  config.UserConfig
	client      user_v1.UserV1Client
	conn        *grpc.ClientConn 
}

func NewUserClient(userConfig config.UserConfig) (client.UserService, error) {
	conn, err := grpc.Dial(userConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	cl := user_v1.NewUserV1Client(conn)
	return &user{
		userConfig: userConfig,
		client:     cl,
		conn:       conn,
	}, nil
}

func (u *user) GetUserAuthData(ctx context.Context, username string) (*user_v1.UserAuthData, error) {
	resp, err := u.client.GetUserAuthData(ctx, &user_v1.GetUserAuthDataRequest{
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (u *user) Close() error {
    if u.conn != nil {
        return u.conn.Close()
    }
    return nil
}
