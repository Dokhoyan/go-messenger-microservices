package user_saver

import (
	"context"

	// "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/repository"
	def "github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
)

const (
	topicName = "user"
)


var _ def.ConsumerService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	consumer       client.KafkaConsumer
}

func NewService(
	userRepository repository.UserRepository,
	consumer client.KafkaConsumer,
) *service {
	return &service{
		userRepository: userRepository,
		consumer:       consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, topicName, s.UserSaveHandler)
	}()

	return errChan
}