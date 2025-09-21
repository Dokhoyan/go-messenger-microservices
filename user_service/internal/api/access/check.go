package access

import (
	"context"
	"errors"

	accesspb "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/access_v1"
	"github.com/golang/protobuf/ptypes/empty"
)

// Check - проверяет доступность ручки для пользователя
func (i *Implementation) Check(ctx context.Context, req *accesspb.CheckRequest) (*empty.Empty, error) {
	//logger.Info("Checking access request", zap.String("endpoint", req.EndpointAddress))

	err := i.service.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, errors.New("access denied")
	}

	return &empty.Empty{}, nil
}