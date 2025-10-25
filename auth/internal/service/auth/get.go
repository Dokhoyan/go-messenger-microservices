package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/utils"
	"github.com/go-redis/redis"
)

func (s *serv) GetAccessToken(ctx context.Context, token string) (string, error) {
	claims, err := utils.VerifyToken(token, s.jwtConfig.RefreshSecretKey())
	if err != nil {
		return "", err
	}

	err = s.checkTokenRefresh(token)
	if err != nil {
		return "", err
	}

	info, err := s.getUserInfoFromStorage(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	log.Printf("claims.Role = %v (%T), info.Role = %v (%T)", claims.Role, claims.Role, info.Role, info.Role)
	if claims.Role != info.Role {
		return "", errors.New("authentication error")
	}

	accessToken, err := utils.GenerateToken(model.UserAuthData{
		Username: claims.Username,
		Role:     claims.Role,
	}, s.jwtConfig.AccessSecretKey(), s.jwtConfig.AccessExpirationTime())
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *serv) GetRefreshToken(ctx context.Context, oldToken string) (string, error) {
	claims, err := utils.VerifyToken(oldToken, s.jwtConfig.RefreshSecretKey())
	if err != nil {
		return "", err
	}

	err = s.checkTokenRefresh(oldToken)
	if err != nil {
		return "", err
	}

	info, err := s.getUserInfoFromStorage(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	if claims.Role != info.Role {
		return "", errors.New("authentication error")
	}

	res := s.redis.Set(oldToken, nil, s.jwtConfig.RefreshExpirationTime())
	if res.Err() != nil {
		return "", err
	}

	refreshToken, err := utils.GenerateToken(model.UserAuthData{
		Username: claims.Username,
		Role:     claims.Role,
	}, s.jwtConfig.RefreshSecretKey(), s.jwtConfig.RefreshExpirationTime())
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *serv) getUserInfoFromStorage(ctx context.Context, username string) (*model.UserAuthData, error) {
	// Пытаемся получить данные из Redis
	res, err := s.redis.Get(username).Result()
	if errors.Is(err, redis.Nil) {
		// Если нет в Redis — достаём из БД
		user, errRep := s.userRepo.Get(ctx, username)
		if errRep != nil {
			return nil, errRep
		}

		return &user.Info, nil
	}
	if err != nil {
		return nil, err
	}

	// Если нашли в Redis — парсим
	var info model.UserAuthData
	if err := json.Unmarshal([]byte(res), &info); err != nil {
		log.Printf("failed to unmarshal Redis user info for %s: %v, raw: %s", username, err, res)
		return nil, err
	}

	log.Printf("Unmarshaled from Redis for %s: %+v", username, info)

	return &info, nil
}


func (s *serv) checkTokenRefresh(refreshToken string) error {
	_, err := s.redis.Get(refreshToken).Result()
	if errors.Is(err, redis.Nil) { 			//refresh-токен валидный, он ещё не использовался и не отозван.
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New("refresh token has expired")		//Если ключ нашёлся в Redis → токен считается просроченным/отозванным
}