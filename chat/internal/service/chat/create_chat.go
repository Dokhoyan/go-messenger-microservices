package chat

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/converter"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
)

// CreateChat creates new chat
func (s *serv) CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (int64, error) {
	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.CreateChat(ctx, req)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logRepository.CreateRecord(ctx,
			converter.ToRecordRepoFromService(id, "create"))
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	s.channels[id] = make(chan *chat_v1.Message, 100)

	return id, nil
}