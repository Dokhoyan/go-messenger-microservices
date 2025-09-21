package user

import (
	"context"
	"github.com/pkg/errors"

	converter "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user/conventer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest)(*desc.CreateResponse, error){
	err:=req.Validate()
	if err!=nil{
		return nil, err
	}
	
	if req.Pass.Password != req.Pass.PasswordConfirm{
		return nil, errors.New("passwords mismatch")
	}

	res, err := i.userservice.Create(ctx, &model.UserCreate{
		Info: converter.ProtoToUserInfo(req.Info),
		Password: req.Pass.Password,
	})

	if err != nil {
		return nil, errors.Errorf("failed to create user: %v", err)
	}

	return &desc.CreateResponse{
		Id: res,
	}, nil
}