package app

import (
	"context"
	"log"

	userImpl "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/closer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	userService "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
)


type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	//dbClient       db.Client
	//txManager      db.TxManager
	pgPool     *pgxpool.Pool

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

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool{
	if s.pgPool==nil{
		pool, err:= pgxpool.Connect(ctx, s.PGConfig().DSN())
		if err!=nil{
			log.Fatalf("failed to connect db: %s",err.Error())
		}

		err=pool.Ping(ctx)
		if err!=nil{
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool=pool
	}

	return s.pgPool
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.PgPool(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userImpl.Implementation {
	if s.userImpl == nil {
		s.userImpl = userImpl.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}