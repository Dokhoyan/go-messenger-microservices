package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/common/pkg/client/db/pg"
	"github.com/Dokhoyan/common/pkg/client/db/transaction"
	"github.com/Dokhoyan/common/pkg/closer"
	"github.com/Dokhoyan/common/pkg/logger"
	"github.com/Dokhoyan/common/pkg/storage"
	cache "github.com/Dokhoyan/common/pkg/storage/redis"
	userImpl "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/auth"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/kafka/producer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	logsRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/logs"
	userRepository "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user"
	accessv1 "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/auth/proto"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	userService "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type serviceProvider struct {
	pgConfig        	 config.PGConfig
	grpcConfig       	 config.GRPCConfig
	httpConfig       	 config.HTTPConfig
	swaggerConfig    	 config.SwaggerConfig
	redisConfig      	 config.RedisConfig
	prometheusConfig 	 config.PrometheusConfig
	kafkaProducerConfig  config.KafkaProducerConfig
	authClientConfig     config.AuthConfig

	redisClient       storage.Redis
	dbClient          db.Client
	txManager         db.TxManager
	kafkaProducer     client.KafkaProducer
	authClient  	  client.Auth


	userRepository   repository.UserRepository
	logsRepo 		 repository.LogsRepository

	userService      service.UserService

	userImpl         *userImpl.Implementation
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

func (s *serviceProvider) AuthClientConfig() config.AuthConfig {
	if s.authClientConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get authClient config: %s", err.Error())
		}

		s.authClientConfig = cfg
	}
	return s.authClientConfig
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

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil{
		cfg, err := config.NewRedisConfig()
		if err != nil{
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.redisConfig = cfg

	}

	return s.redisConfig

}

func (s *serviceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil{
		cfg, err := config.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka config: %s", err.Error())
		}

		s.kafkaProducerConfig = cfg

	}
	return s.kafkaProducerConfig
}

func (s *serviceProvider) AuthClient(ctx context.Context) client.Auth {
	if s.authClient == nil {
		conn, err := grpc.DialContext(ctx, s.AuthClientConfig().Address(), grpc.WithTransportCredentials(insecure.NewCredentials()),)
		if err != nil {
			logger.Fatal("failed to connect localhost:50053", zap.Error(err))
		}
		closer.Add(conn.Close)

		client := accessv1.NewAccessV1Client(conn)
		s.authClient = auth.NewClient(client)
	}

	return s.authClient
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

func (s *serviceProvider) KafkaProducer() client.KafkaProducer{
	if s.kafkaProducer == nil {
		p, err := producer.NewProducer(s.KafkaProducerConfig().Brokers())
		if err != nil {
			log.Fatalf("failed new kafka producer: %v", err)
		}
		s.kafkaProducer = p

		closer.Add(s.kafkaProducer.Close)
	}

	return s.kafkaProducer
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
			s.RedisClient(),
			s.KafkaProducer(),)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userImpl.Implementation {
	if s.userImpl == nil {
		s.userImpl = userImpl.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

