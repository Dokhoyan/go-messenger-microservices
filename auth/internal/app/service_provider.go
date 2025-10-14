package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Dokhoyan/common/pkg/closer"
	"github.com/Dokhoyan/common/pkg/storage"
	cache "github.com/Dokhoyan/common/pkg/storage/redis"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/api/access"
	authApi "github.com/Dokhoyan/go-messenger-microservices/auth/internal/api/auth"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client"
	userClient "github.com/Dokhoyan/go-messenger-microservices/auth/internal/client/user"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service/access"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service/auth"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/utils"
	"github.com/go-redis/redis"
)

type serviceProvider struct {
	grpcConfig       config.GRPCConfig
	redisConfig      config.RedisConfig
	jwtConfig        config.JWTConfig
	userConfig       config.UserConfig

	redisClient      storage.Redis
	userClient       client.UserService

	
	authService      service.AuthService
	accessService    service.AccessService

	authImpl         *authApi.Implementation
	accessImpl       *accessApi.Implementation

	accessChecker    utils.AccessChecker
}

func newServiceProvider() (*serviceProvider) {
	return &serviceProvider{}
}


func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %v", err)
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig{
	if s.redisConfig == nil{
		cfg, err:=config.NewRedisConfig()
		if err!=nil{
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.redisConfig=cfg

	}

	return s.redisConfig

}

func (s *serviceProvider) UserConfig() config.UserConfig{
	if s.userConfig == nil{
		cfg, err:=config.NewUserConfig()
		if err!=nil{
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.userConfig=cfg
	}

	return s.userConfig
}

func (s *serviceProvider) UserClient() client.UserService{
	if s.userClient == nil{
		cl , err := userClient.NewUserClient(s.UserConfig())
		if err != nil{
			log.Fatal("error conn UserClient")
		}

		closer.Add(cl.Close)

		s.userClient = cl
	}

	return s.userClient
}

func (s *serviceProvider) RedisClient() storage.Redis {
	if s.redisClient==nil{
		cl, err := cache.NewRedisConnection(&redis.Options{
			Addr:      s.RedisConfig().Address(),
			Password:  s.RedisConfig().Password(),
			DB:        0,
		})

		if err != nil {
			log.Fatalf("failed to create redis client: %v", err)
		}

		err = cl.Ping()
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}

		closer.Add(cl.Close)

		s.redisClient = cl

		s.routesMigrate()
	}

	return s.redisClient
}

//миграция "прав доступа" в Redis.
func (s *serviceProvider) routesMigrate(){
	routes := s.redisConfig.RoutesAccesses()

	for route, roles := range routes{
		roleJSON, err := json.Marshal(roles)
		if err!=nil{
			log.Fatalf("error at json marshal")
		}

		_, err = s.RedisClient().Set(route, roleJSON, 0).Result()
		if err!=nil{
			log.Fatalf("error at migration routes to redis")
		}
	}
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = access.NewService(
			s.JWTConfig(),
			s.AccessChecker(ctx))
	}

	return s.accessService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = auth.NewService(
			s.RedisClient(),
			s.UserClient(),
			s.JWTConfig(),
		)
	}

	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *authApi.Implementation {
	if s.authImpl == nil {
		s.authImpl = authApi.NewImplementation(
			s.AuthService(ctx),
		)
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *accessApi.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = accessApi.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) AccessChecker(ctx context.Context) utils.AccessChecker {
	if s.accessChecker == nil {
		s.accessChecker = utils.NewRouteAccessChecker(
			s.JWTConfig(),
			s.RedisClient(),
			s.UserClient(),
		)
	}

	return s.accessChecker
}