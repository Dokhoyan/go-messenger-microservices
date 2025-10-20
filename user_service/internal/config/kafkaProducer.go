package config

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

type kafkaProducerConfig struct {
	brokers []string
}

func NewKafkaProducerConfig() (KafkaProducerConfig, error) {
	brokerEnv := os.Getenv("KAFKA_BROKERS")
	if brokerEnv == "" {
		return nil, errors.New("KAFKA_BROKERS not found")
	}

	return &kafkaProducerConfig{
		brokers: strings.Split(brokerEnv, ","),
	}, nil
}

func (c *kafkaProducerConfig) Brokers() []string {
	return c.brokers
}

func (c *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	return config
}
