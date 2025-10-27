package chat

import (
	"context"
)

// GetChatIDByName calls repo for get chat id by name
func (s *serv) GetChatIDByName(ctx context.Context, chatname string) (int64, error) {
	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.GetChatIDByName(ctx, chatname)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}