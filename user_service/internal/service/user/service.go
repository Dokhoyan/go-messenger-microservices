package user

import (
	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/common/pkg/storage"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
)

type serv struct{
	userRepository repository.UserRepository
	txManager      db.TxManager
	logsRepo       repository.LogsRepository
	storage   	   storage.Redis
	producer       client.KafkaProducer     
}

func NewService (userRepository repository.UserRepository, txManager db.TxManager, logsRepo repository.LogsRepository, storage storage.Redis, producer client.KafkaProducer) service.UserService{
	return &serv{userRepository: userRepository,
				 txManager: txManager,
				 logsRepo:  logsRepo,
				 storage: storage,
				 producer: producer,
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