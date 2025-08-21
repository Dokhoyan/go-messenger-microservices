package user

import (
	//"context"

	//"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/converter"
	//"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
)

type serv struct{
	userRepository repository.UserRepository
	
}

func NewService (userRepository repository.UserRepository) service.UserService{
	return &serv{userRepository: userRepository,
	}
}

