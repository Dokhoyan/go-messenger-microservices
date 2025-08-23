package converter

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	modelRepo "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user/model"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:   info.Name,
		Username: info.Username,
		Email:   info.Email,
		Birth_date: info.Birth_date,
		Avatar_url: info.Avatar_url,
	}
}



