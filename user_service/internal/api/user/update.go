package user

import (
	"context"

	converter "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user/conventer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error){
	err := i.userservice.Update(ctx, &model.UserUpdate{
		ID: req.Id,
		Info: converter.ProtoToUserInfoUpdate(req.GetInfo()),
	})

	if err != nil {
		return nil, errors.Errorf("failed to update user: %v", err)
	}

	return &empty.Empty{}, nil
}
//gRPC не позволяет возвращать nil — у каждого RPC-метода должен быть объект ответа.
