package user

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserInfo) (int64, error){
	id, err:=s.userRepository.Create(ctx,info)
	if err!=nil{
		return 0, err
	}

	return id, nil
}