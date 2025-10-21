package interceptor

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type authInterceptor struct {
	authClient client.Auth
}

func NewAuthInterceptor(authClient client.Auth) *authInterceptor {
	return &authInterceptor{
		authClient: authClient,
	}
}

func (i *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		ok, err = i.authClient.Check(ctx, info.FullMethod)
		if err != nil || !ok {
			return nil, err
		}

		return handler(ctx, req)
	}
}