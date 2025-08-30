package user

import (
	"context"
	"log"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/converter"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

func(i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error){
	userObj,err:=i.userservice.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, username: %s, email: %s, birth_date: %v, avatar_url: %s\n",
	 userObj.ID, userObj.Info.Name, userObj.Info.Username, 
	userObj.Info.Email, userObj.Info.Birth_date, userObj.Info.Avatar_url)

	return &desc.GetResponse{
		User: converter.ToUserFromService(userObj),
	}, nil
}

