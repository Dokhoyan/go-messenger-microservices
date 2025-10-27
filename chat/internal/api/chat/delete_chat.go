package chat

import (
	"context"
	"log"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteChat(ctx context.Context, req *chat_v1.DeleteChatRequest) (*emptypb.Empty, error) {
	err := i.chatService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("deleted chat with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}