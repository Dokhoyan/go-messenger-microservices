package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/utils"
)

// берет юзера из базы по юзернейму, сравнивает пароли, генерит рефреш токен и сохраняет инфо юзера в кэш
func (s *serv) Login(ctx context.Context, req model.LoginDTO) (string, error) {
	

	user, err := s.userRepo.Get(ctx, req.Username)
	if err != nil {
		return "", err
	}


    if !utils.VerifyPassword(user.Info.PasswordHash, req.Password) {
		return "", errors.New("authentication failed passw")
	}


	token, err := utils.GenerateToken(user.Info, s.jwtConfig.RefreshSecretKey(), s.jwtConfig.RefreshExpirationTime())
	if err != nil {
		return "", err
	}

	infoJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	res := s.redis.Set(user.Info.Username, infoJSON, 0)
	if res.Err() != nil {
		return "", err
	}

	return token, nil
}