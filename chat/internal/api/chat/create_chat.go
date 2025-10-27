package chat

import (
	"context"
	"log"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
)

func (i *Implementation) CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (*chat_v1.CreateChatResponse, error) {
	id, err := i.chatService.CreateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted chat with id: %d", id)
	
	return &chat_v1.CreateChatResponse{
		Id : id,
		}, nil
}