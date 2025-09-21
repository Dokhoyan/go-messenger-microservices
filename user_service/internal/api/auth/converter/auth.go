package converter

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/auth_v1"
)

func AuthProtoToAuthDTO(req *desc.LoginRequest) model.LoginDTO {
	return model.LoginDTO{
		Password: req.GetPassword(),
		Username: req.GetUsername(),
	}
}