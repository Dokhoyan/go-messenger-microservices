package converter

import (
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/model"
)

func ToRecordRepoFromService(chatID int64, action string) *model.Record {
	return &model.Record{
		ChatID:    chatID,
		CreatedAt: time.Now(),
		Action:    action,
	}
}