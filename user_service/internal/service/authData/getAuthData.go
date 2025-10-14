package authdata

import (
	"context"

	"github.com/Dokhoyan/common/pkg/filter"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
)

func (s *serv) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user *model.User

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		conditions := filter.MakeFilter(filter.Condition{
			Key:   model.UserNameFieldCode,
			Value: username,
		})

		user, errTx = s.userRepo.Get(ctx, conditions)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logsRepo.Create(ctx, model.Log{
			Action:  "user login",
			Content: username,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}


