package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service/auth/converter"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/utils"
)

// берет юзера из базы по юзернейму, сравнивает пароли, генерит рефреш токен и сохраняет инфо юзера в кэш
func (s *serv) Login(ctx context.Context, req model.LoginDTO) (string, error) {
	

	user, err := s.userClient.GetUserAuthData(ctx, req.Username)
	if err != nil {
		return "", err
	}

	userAuthData := converter.ProtoToUser(user)

    if !utils.VerifyPassword(userAuthData.Password, req.Password) {
		return "", errors.New("authentication failed passw")
	}

	token, err := utils.GenerateToken(*userAuthData, s.jwtConfig.RefreshSecretKey(), s.jwtConfig.RefreshExpirationTime())
	if err != nil {
		return "", err
	}

	infoJSON, err := json.Marshal(userAuthData)
	if err != nil {
		return "", err
	}
	res := s.redis.Set(userAuthData.Username, infoJSON, 0)
	if res.Err() != nil {
		return "", err
	}

	return token, nil
}