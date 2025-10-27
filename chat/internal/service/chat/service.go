package chat

import (
	"sync"

	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
)

type Chat struct {
	userConnections map[string]chat_v1.ChatV1_ConnectChatServer
	m               sync.RWMutex
}

type serv struct {
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager

	chats  map[int64]*Chat
	mxChat sync.RWMutex

	channels  map[int64]chan *chat_v1.Message
	mxChannel sync.RWMutex
}

func NewService(
	chatRepository repository.ChatRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		logRepository:  logRepository,
		txManager:      txManager,
		chats:          make(map[int64]*Chat),
		channels:       make(map[int64]chan *chat_v1.Message),
	}
}