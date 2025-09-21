package user

import (
	"context"

	"github.com/Dokhoyan/common/pkg/filter"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		conditions := filter.MakeFilter(filter.Condition{
			Key:   model.IDFieldCode,
			Value: id,
		})

		user, errTx = s.userRepository.Get(ctx, conditions)
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