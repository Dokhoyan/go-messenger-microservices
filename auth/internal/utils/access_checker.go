package utils

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Dokhoyan/common/pkg/storage"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"

	//"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/go-redis/redis"
)

type AccessChecker interface {
	AccessCheck(ctx context.Context, token string, endpoint string) (bool, error)
}

type routeAccessChecker struct {
	jwtConfig config.JWTConfig
	redis     storage.Redis
	userClient client.UserService
}

// NewRouteAccessChecker - создает новый экземпляр верификатора
func NewRouteAccessChecker(jwtCfg config.JWTConfig, redis storage.Redis, userClient client.UserService) AccessChecker {
	return &routeAccessChecker{
		jwtConfig: jwtCfg,
		redis:     redis,
		userClient: userClient,
	}
}

func (r *routeAccessChecker) AccessCheck(ctx context.Context, token string, endpoint string) (bool, error){
	claims, err := VerifyToken(token, r.jwtConfig.AccessSecretKey())
	if err!=nil{
		return false, err
	}

	var info *model.UserAuthData

	res, err := r.redis.Get(claims.Username).Result()
	if errors.Is(err, redis.Nil) {
		

		user, errRep := r.userClient.GetUserAuthData(ctx, claims.Username)
		if errRep != nil{
			return false, err
		}

		info = &model.UserAuthData{
			Username: user.Username,
			Role: model.UserRole(user.Role),
			Password: user.PasswordHash,
		}
	}

	if err != nil {
		return false, err
	}
	
	if info == nil{
		err=json.Unmarshal([]byte(res), &info)
		if err!=nil{
			return false, err
		}
	}

	if info.Role == model.ADMIN && claims.Role == model.ADMIN {
		return true, nil
	}

	res, err = r.redis.Get(endpoint).Result()
	if errors.Is(err, redis.Nil) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	var roles []model.UserRole
	err = json.Unmarshal([]byte(res), &roles)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role == info.Role {
			return true, nil
		}
	}

	return false, nil

	
}