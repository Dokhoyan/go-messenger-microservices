package user

import (
	"context"
	"strconv"
	commonVal "github.com/Dokhoyan/common/pkg/sys/validate"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/validate"
)

func (s *serv) Delete(ctx context.Context, id int64) (error){
	err := commonVal.Validate(
		ctx,
		validate.ID(id),
	)
	if err != nil {
		return err
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error{
		var errTx error

		errTx = s.userRepository.Delete(ctx, id)
		if errTx != nil{
			return errTx
		}

		_, errTx = s.logsRepo.Create(ctx, model.Log{
			Action:  "user deleted",
			Content: strconv.FormatInt(id, 10),
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