package chat

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/converter"
)

// DeleteChat removes existing chat
func (s *serv) DeleteChat(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.chatRepository.DeleteChat(ctx, id)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logRepository.CreateRecord(ctx,
			converter.ToRecordRepoFromService(id, "delete"))
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}