package user

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

type Implementation struct{
	desc.UnimplementedUserV1Server
	userservice service.UserService
}

func NewImplementation(userservice service.UserService) *Implementation {
	return &Implementation{
		userservice: userservice,
	}
}


