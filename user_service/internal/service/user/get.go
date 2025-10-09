package user

import (
	"context"
	"strconv"

	"github.com/Dokhoyan/common/pkg/filter"
	commonVal "github.com/Dokhoyan/common/pkg/sys/validate"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/validate"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	err := commonVal.Validate(
		ctx,
		validate.ID(id),
	)
	if err != nil{
		return nil, err
	}
	
	var user *model.User

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		conditions := filter.MakeFilter(filter.Condition{
			Key:   model.IDFieldCode,
			Value: id,
		})

		user, errTx = s.userRepository.Get(ctx, conditions)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logsRepo.Create(ctx, model.Log{
			Action:  "user fetch",
			Content: strconv.FormatInt(id, 10),
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