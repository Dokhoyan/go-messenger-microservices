package repository

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
)

type ChatRepository interface {
	GetChatIDByName(ctx context.Context, chatname string) (int64, error)
	CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
}

type LogRepository interface {
	CreateRecord(ctx context.Context, record *model.Record) (int64, error)
}