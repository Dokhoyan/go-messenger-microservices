package authdata

import (
	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
)

type serv struct{
	userRepo repository.UserRepository
	txManager      db.TxManager
	logsRepo       repository.LogsRepository
}

func NewService (userRepository repository.UserRepository, txManager db.TxManager, logsRepo repository.LogsRepository) service.AuthDataService{
	return &serv{userRepo: userRepository,
				 txManager: txManager,
				 logsRepo:  logsRepo,

	}
}
