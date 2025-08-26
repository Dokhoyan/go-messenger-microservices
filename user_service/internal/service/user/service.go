package user

import (
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
)

type serv struct{
	userRepository repository.UserRepository
	txManager      db.TxManager
	
}

func NewService (userRepository repository.UserRepository, txManager db.TxManager) service.UserService{
	return &serv{userRepository: userRepository,
				 txManager: txManager,
	}
}

