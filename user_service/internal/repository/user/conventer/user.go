package converter

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	modelRepo "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user/model"
)

func ToUserFromRepo(note *modelRepo.User) *model.User {
	return &model.User{
		ID:        note.ID,
		Info:      ToUserInfoFromRepo(note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
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



