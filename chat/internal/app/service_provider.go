package app

import (
	"context"
	"log"

	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/common/pkg/client/db/pg"
	"github.com/Dokhoyan/common/pkg/client/db/transaction"
	"github.com/Dokhoyan/common/pkg/closer"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/api/chat"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/client/auth"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/interceptor"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/repository"
	chatRepo "github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/repository/chat"
	logRepo "github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/repository/log"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/service"
	chatServ "github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/service/chat"
)

type serviceProvider struct {
	pgConfig       config.PGConfig
	grpcConfig     config.GRPCConfig
	grpcAuthConfig config.AuthConfig

	authClient client.AuthService
	dbClient       db.Client
	txManager      db.TxManager

	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository

	chatService    service.ChatService

	chatImpl *chat.Implementation

	accessChecker interceptor.AccessChecker
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig returns new PGConfig
func (s *serviceProvider) PGConfig() (config.PGConfig, error) {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			return nil, err
		}

		s.pgConfig = cfg
	}

	return s.pgConfig, nil
}

// GRPCConfig returns new GRPCConfig
func (s *serviceProvider) GRPCConfig() (config.GRPCConfig, error) {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			return nil, err
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig, nil
}

// GRPCAuthConfig returns new GRPCAuthConfig
func (s *serviceProvider) AuthConfig() (config.AuthConfig) {
	if s.grpcAuthConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to load auth config: %v", err)
		}

		s.grpcAuthConfig = cfg
	}

	return s.grpcAuthConfig
}


// DBClient returns new db client
func (s *serviceProvider) DBClient(ctx context.Context) (db.Client, error) {
	if s.dbClient == nil {
		pgConfig, err := s.PGConfig()
		if err != nil {
			return nil, err
		}
		cl, err := pg.New(ctx, pgConfig.DSN())
		if err != nil {
			return nil, err
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			return nil, err
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient, nil
}

// TxManager returns new db TxManager
func (s *serviceProvider) TxManager(ctx context.Context) (db.TxManager, error) {
	if s.txManager == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.txManager = transaction.NewTransactionManager(dbClient.DB())
	}

	return s.txManager, nil
}

// ChatRepository returns new ChatRepository
func (s *serviceProvider) ChatRepository(ctx context.Context) (repository.ChatRepository, error) {
	if s.chatRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.chatRepository = chatRepo.NewRepository(dbClient)
	}

	return s.chatRepository, nil
}

// LogRepository returns new LogRepository
func (s *serviceProvider) LogRepository(ctx context.Context) (repository.LogRepository, error) {
	if s.logRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.logRepository = logRepo.NewRepository(dbClient)
	}

	return s.logRepository, nil
}

// ChatService returns new ChatService
func (s *serviceProvider) ChatService(ctx context.Context) (service.ChatService) {
	if s.chatService == nil {
		chatRepo, err := s.ChatRepository(ctx)
		if err != nil {
			log.Fatalf("failed to create chat repository: %v", err)
		}
		logRepo, err := s.LogRepository(ctx)
		if err != nil {
			log.Fatalf("failed to create log repository: %v", err)
		}
		txManager, err := s.TxManager(ctx)
		if err != nil {
			log.Fatalf("failed to create tx manager: %v", err)
		}
		s.chatService = chatServ.NewService(
			chatRepo, logRepo, txManager,
		)
	}

	return s.chatService
}

// ChatImpl returns new Chat Service implementation
func (s *serviceProvider) ChatImplementation(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.ChatService(ctx))
	}

	return s.chatImpl
}

func (s *serviceProvider) AuthClient(_ context.Context) client.AuthService {
	if s.authClient == nil {
		authClient, err := auth.NewAuthClient(s.AuthConfig())
		if err != nil {
			log.Fatalf("failed to connect auth service: %v", err)
		}

		s.authClient = authClient
	}

	return s.authClient
}

func (s *serviceProvider) AccessChecker(ctx context.Context) interceptor.AccessChecker {
	if s.accessChecker == nil {
		s.accessChecker = interceptor.NewAccessChecker(s.AuthClient(ctx))
	}

	return s.accessChecker
}