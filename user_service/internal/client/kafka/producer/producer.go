package producer

import (
	"context"
	"encoding/json"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/IBM/sarama"
)

var _ client.KafkaProducer = (*producer)(nil)


type producer struct {
	syncProducer sarama.SyncProducer
	brokers      []string
}

// NewProducer — создаёт новый Kafka-продюсер
func NewProducer(brokers []string) (*producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &producer{
		syncProducer: p,
		brokers:      brokers,
	}, nil
}

func (p *producer) Produce(ctx context.Context, topicName string, handler client.KafkaHandler) error {
	data, err := handler.Data()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = p.syncProducer.SendMessage(msg)
	return err
}

// Close — закрывает соединение с Kafka
func (p *producer) Close() error {
	if p.syncProducer != nil {
		return p.syncProducer.Close()
	}
	return nil
}

type UserCreatedHandler struct {
	Username     string
	PasswordHash string
	Role         model.UserRole
}


func (h *UserCreatedHandler) Data() (interface{}, error) {
	return map[string]interface{}{
		"username":     h.Username,
		"passwordHash": h.PasswordHash, 
		"role":         h.Role, 
	}, nil
}