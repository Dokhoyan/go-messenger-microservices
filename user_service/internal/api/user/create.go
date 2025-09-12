package user

import (
	"context"
	"log"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/converter"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest)(*desc.CreateResponse, error){
	err:=req.Validate()
	if err!=nil{
		return nil, err
	}
	
	id, err:=i.userservice.Create(ctx, converter.ToUserInfoFromDesc(req.GetInfo()))
	if err!=nil{
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)
	
	return &desc.CreateResponse{
		Id: id,
	}, nil
}