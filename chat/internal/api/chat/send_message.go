package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage sends message to server
func (i *Implementation) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (
	*emptypb.Empty, error,
) {
	log.Printf("Send message from %s with text: %s",
		req.GetMessage().GetFrom(), req.GetMessage().GetText())

	chatID, err := i.chatService.GetChatIDByName(ctx, req.Chatname)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = i.chatService.SendMessage(ctx, chatID, req.Message)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("error while sending message %v", err.Error())
	}

	return &emptypb.Empty{}, nil
}
