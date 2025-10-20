package user

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

type Implementation struct{
	desc.UnimplementedUserV1Server
	userservice service.UserService
	authDataService service.AuthDataService
}

func NewImplementation(userservice service.UserService, authDataService service.AuthDataService) *Implementation {
	return &Implementation{
		userservice: userservice,
		authDataService: authDataService,
	}
}


