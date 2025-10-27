package service

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
)

type ChatService interface {
	ConnectChat(ctx context.Context, chatID int64, username string, stream chat_v1.ChatV1_ConnectChatServer) error
	SendMessage(ctx context.Context, chatID int64, message *chat_v1.Message) error
	CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	GetChatIDByName(ctx context.Context, chatname string) (int64, error)
}