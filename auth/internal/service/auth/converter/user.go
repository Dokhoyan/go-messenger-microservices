package converter

import (
	user_v1 "github.com/Dokhoyan/go-messenger-microservices/auth/internal/client/user/proto"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
)

func ProtoToUser(user *user_v1.UserAuthData) *model.UserAuthData {
	return &model.UserAuthData{
		Username: user.Username,
		Role:     model.UserRole(user.Role),
		Password: user.PasswordHash,
		
	}
}
