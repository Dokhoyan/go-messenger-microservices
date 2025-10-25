package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
	"github.com/IBM/sarama"
)

func (s *service) NoteSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {

	log.Printf("Received message: %s", string(msg.Value))
	
	userInfo := &model.UserAuthData{}
	err := json.Unmarshal(msg.Value, userInfo)
	if err != nil {
		log.Printf("Raw message: %s", string(msg.Value))
		return err
	}

	id, err := s.userRepository.Create(ctx, userInfo)
	log.Printf("Failed to insert user into auth DB: %v", err)
	if err != nil {
		return err
	}

	log.Printf("User with id %d created\n", id)

	return nil
}