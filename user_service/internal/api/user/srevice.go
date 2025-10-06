package user

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

type Implementation struct{
	desc.UnimplementedUserV1Server
	userservice service.UserService
	accessService service.AccessService
}

func NewImplementation(userservice service.UserService, access service.AccessService) *Implementation {
	return &Implementation{
		userservice: userservice,
		accessService: access,
	}
}