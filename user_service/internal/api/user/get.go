package user

import (
	"context"

	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

func(s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error){

	err:=req.Validate()   //прото валидация
	if err!=nil{
		return nil, err
	}


	return nil ,nil
}