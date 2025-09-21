package user

import (
	"context"
	

	converter "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user/conventer"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"

	
	"github.com/pkg/errors"
)

func(i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error){
	userObj, err := i.userservice.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.Errorf("failed to get user: %v", err)
	}

	return &desc.GetResponse{
		User: converter.UserToProto(userObj),
	}, nil
}

