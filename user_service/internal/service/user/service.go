package user

import (
	"github.com/Dokhoyan/common/pkg/storage"
	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
)

type serv struct{
	userRepository repository.UserRepository
	txManager      db.TxManager
	logsRepo       repository.LogsRepository
	storage   	   storage.Redis
	
}

func NewService (userRepository repository.UserRepository, txManager db.TxManager, logsRepo repository.LogsRepository, storage storage.Redis ) service.UserService{
	return &serv{userRepository: userRepository,
				 txManager: txManager,
				 logsRepo:  logsRepo,
				 storage: storage,
	}
}

func NewMockService(deps ...interface{}) service.UserService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			srv.userRepository = s
		}
	}

	return &srv
}