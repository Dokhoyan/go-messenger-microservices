package authservice

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Dokhoyan/common/pkg/filter"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/utils"
)

// берет юзера из базы по юзернейму, сравнивает пароли, генерит рефреш токен и сохраняет инфо юзера в кэш
func (s *serv) Login(ctx context.Context, req model.LoginDTO) (string, error) {
	conditions:=filter.MakeFilter(filter.Condition{
		Key: model.UserNameFieldCode,
		Value: req.Password,                   //возможно ошибка 
	})

	user, err := s.userRepo.Get(ctx, conditions)
	if err != nil {
		return "", err
	}

    if !utils.VerifyPassword(user.Password, req.Password) {
		return "", errors.New("authentication failed")
	}

	token, err := utils.GenerateToken(user.Info, s.jwtConfig.RefreshSecretKey(), s.jwtConfig.RefreshExpirationTime())
	if err != nil {
		return "", err
	}

	infoJSON, err := json.Marshal(user.Info)
	if err != nil {
		return "", err
	}
	res := s.redis.Set(user.Info.Username, infoJSON, 0)
	if res.Err() != nil {
		return "", err
	}

	return token, nil
}