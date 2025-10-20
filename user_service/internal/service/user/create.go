package user

import (
	"context"
	"strconv"

	"github.com/Dokhoyan/common/pkg/filter"
	commonVal "github.com/Dokhoyan/common/pkg/sys/validate"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/kafka/producer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/utils"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/validate"
	"github.com/pkg/errors"
)

const (
	topicName     = "user"
)

func (s *serv) Create(ctx context.Context, userParams *model.UserCreate) (int64, error) {
	err := commonVal.Validate(
		ctx,
		validate.Email(userParams.Info.Email),
		validate.Birthday(userParams.Info.Birth_date),
	)
	if err != nil {
		return 0, err
	}

	conditions := filter.MakeFilter(filter.Condition{
		Key:   model.UserNameFieldCode,
		Value: userParams.Info.Username,
	})

	user, err := s.userRepository.Get(ctx, conditions)
	if err != nil && err.Error() != "user not found" {
		return 0, err
	}

	if user != nil {
		return 0, errors.Errorf(`user with username "%s" already exist`, userParams.Info.Username)
	}

	hashedPassword, err := utils.HashPassword(userParams.Password)
	if err != nil {
		return 0, errors.Errorf("failed hash password: %v", err)
	}

	userParams.Password = hashedPassword

	var id int64

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		id, errTx = s.userRepository.Create(ctx, userParams)
		if errTx != nil {
			return errTx
		}

		handler := &producer.UserCreatedHandler{
			Username: userParams.Info.Username,
			PasswordHash: userParams.Password,
			Role: model.UserRole.String(userParams.Info.Role),
		}
		errTx = s.producer.Produce(ctx, topicName, handler )
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logsRepo.Create(ctx, model.Log{
			Action:  "user created",
			Content: strconv.FormatInt(id, 10),
		})
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