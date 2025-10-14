package converter

import (
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
	desc "github.com/Dokhoyan/go-messenger-microservices/auth/pkg/api/auth_v1"
)

func AuthProtoToAuthDTO(req *desc.LoginRequest) model.LoginDTO {
	return model.LoginDTO{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
}