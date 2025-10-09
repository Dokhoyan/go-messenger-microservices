package user

import (
	"context"
	"strconv"

	commonVal "github.com/Dokhoyan/common/pkg/sys/validate"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/validate"
)

func (s *serv) Update(ctx context.Context, user *model.UserUpdate) error {
	err := commonVal.Validate(
		ctx,
		validate.Email(user.Info.Email),
		validate.Birthday(user.Info.Birth_date),
	)
	if err != nil {
		return err
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.userRepository.Update(ctx, user)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logsRepo.Create(ctx, model.Log{
			Action:  "user updated",
			Content: strconv.FormatInt(user.ID, 10),
		})
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