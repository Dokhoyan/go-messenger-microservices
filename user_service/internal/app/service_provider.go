package app

import (
	"context"
	"log"

	userImpl "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/db/pg"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/db/transaction"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/closer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	userService "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
)


type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient       db.Client
	txManager      db.TxManager


	userRepository repository.UserRepository

	userService service.UserService

	userImpl *userImpl.Implementation
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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
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

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.TxManager(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userImpl.Implementation {
	if s.userImpl == nil {
		s.userImpl = userImpl.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}