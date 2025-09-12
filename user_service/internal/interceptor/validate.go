package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

type Validator interface {
	Validate() error
}


func ValidateInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if val, ok := req.(Validator); ok {
		err := val.Validate()
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}