package chat

import (
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
)

type Implementation struct {
	chat_v1.UnimplementedChatV1Server
	chatService service.ChatService
}

func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}