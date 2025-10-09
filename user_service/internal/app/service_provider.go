package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Dokhoyan/common/pkg/storage"
	cache "github.com/Dokhoyan/common/pkg/storage/redis"
	accessAPI "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/access"
	accessImpl "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/access"
	authAPI "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/auth"
	authImpl "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/auth"
	userImpl "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/common/pkg/client/db/pg"
	"github.com/Dokhoyan/common/pkg/client/db/transaction"
	"github.com/Dokhoyan/common/pkg/closer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	logsRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/logs"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	accessServ "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/access"
	authservice "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/auth"
	userService "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/utils"
	"github.com/go-redis/redis"
)


type serviceProvider struct {
	pgConfig         config.PGConfig
	grpcConfig       config.GRPCConfig
	httpConfig       config.HTTPConfig
	swaggerConfig    config.SwaggerConfig
	redisConfig      config.RedisConfig
	jwtConfig        config.JWTConfig
	prometheusConfig config.PrometheusConfig

	redisClient      storage.Redis
	dbClient         db.Client
	txManager        db.TxManager


	userRepository   repository.UserRepository
	logsRepo 		 repository.LogsRepository

	userService      service.UserService
	authService      service.AuthService
	accessService    service.AccessService

	userImpl         *userImpl.Implementation
	authImpl         *authImpl.Implementation
	accessImpl       *accessImpl.Implementation

	accessChecker    utils.AccessChecker
}

func newServiceProvider() (*serviceProvider) {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig==nil{
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig=cfg
	}
	return s.pgConfig
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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) PrometheusConfig() config.PrometheusConfig {
	if s.prometheusConfig == nil {
		cfg, err := config.NewPrometheusConfig()
		if err != nil {
			log.Fatalf("failed to get prometheus config: %v", err)
		}

		s.prometheusConfig = cfg
	}

	return s.prometheusConfig
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

func(s *serviceProvider) RedisConfig() config.RedisConfig{
	if s.redisConfig == nil{
		cfg, err:=config.NewRedisConfig()
		if err!=nil{
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.redisConfig=cfg

	}

	return s.redisConfig

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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) LogsRepository(ctx context.Context) repository.LogsRepository {
	if s.logsRepo == nil {
		s.logsRepo = logsRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logsRepo
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
		    s.TxManager(ctx), 
			s.LogsRepository(ctx), 
			s.RedisClient())
	}

	return s.userService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessServ.NewService(
			s.JWTConfig(),
			s.AccessChecker(ctx))
	}

	return s.accessService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authservice.NewService(
			s.RedisClient(),
			s.UserRepository(ctx),
			s.JWTConfig(),
		)
	}

	return s.authService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userImpl.Implementation {
	if s.userImpl == nil {
		s.userImpl = userImpl.NewImplementation(s.UserService(ctx), s.AccessService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *authAPI.Implementation {
	if s.authImpl == nil {
		s.authImpl = authAPI.NewImplementation(
			s.AuthService(ctx),
		)
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *accessAPI.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = accessAPI.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) AccessChecker(ctx context.Context) utils.AccessChecker {
	if s.accessChecker == nil {
		s.accessChecker = utils.NewRouteAccessChecker(
			s.JWTConfig(),
			s.RedisClient(),
			s.UserRepository(ctx),
		)
	}

	return s.accessChecker
}