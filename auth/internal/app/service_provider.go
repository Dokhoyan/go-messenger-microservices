package app

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/common/pkg/client/db/pg"
	"github.com/Dokhoyan/common/pkg/closer"
	"github.com/Dokhoyan/common/pkg/storage"
	cache "github.com/Dokhoyan/common/pkg/storage/redis"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/api/access"
	authApi "github.com/Dokhoyan/go-messenger-microservices/auth/internal/api/auth"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client/kafka/consumer"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/repository/user"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service/access"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service/auth"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service/consumer/user_saver"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/utils"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis"
)

func init() {
    sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
}
type serviceProvider struct {
	grpcConfig          config.GRPCConfig
	pgConfig            config.PGConfig
	redisConfig         config.RedisConfig
	jwtConfig           config.JWTConfig
	kafkaConsumerConfig config.KafkaConsumerConfig
	
	dbClient 			db.Client
	redisClient         storage.Redis

	userRepository 		repository.UserRepository

	userSaverConsumer 	service.ConsumerService
	authService       	service.AuthService
	accessService     	service.AccessService

	authImpl         	*authApi.Implementation
	accessImpl       	*accessApi.Implementation

	accessChecker    	utils.AccessChecker

	consumer             client.KafkaConsumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *consumer.GroupHandler
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

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := config.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
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
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = user.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = user_saver.NewService(
			s.UserRepository(ctx),
			s.Consumer(),
		)
	}

	return s.userSaverConsumer
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
			s.UserRepository(ctx),
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
			s.UserRepository(ctx),
		)
	}

	return s.accessChecker
}

func (s *serviceProvider) Consumer() client.KafkaConsumer {
	if s.consumer == nil {
		s.consumer = consumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *consumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = consumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}