package accessApi

import (
	"context"
	"errors"

	accesspb "github.com/Dokhoyan/go-messenger-microservices/auth/pkg/api/access_v1"
	"github.com/golang/protobuf/ptypes/empty"
)

//проверяет доступность ручки для пользователя
func (i *Implementation) Check(ctx context.Context, req *accesspb.CheckRequest) (*empty.Empty, error) {

	err := i.service.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, errors.New("access denied")
	}

	return &empty.Empty{}, nil
}