package chat

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConnectChat connects to chat
func (s *serv) ConnectChat(_ context.Context, chatID int64, username string, stream chat_v1.ChatV1_ConnectChatServer) error {
	s.mxChannel.RLock()
	chatChan, ok := s.channels[chatID]
	s.mxChannel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	s.mxChat.Lock()
	// Create new chat in map if it doesnt exist
	if _, okChat := s.chats[chatID]; !okChat {
		s.chats[chatID] = &Chat{
			userConnections: make(map[string]chat_v1.ChatV1_ConnectChatServer),
		}
	}
	s.mxChat.Unlock()

	s.chats[chatID].m.Lock()
	s.chats[chatID].userConnections[username] = stream
	s.chats[chatID].m.Unlock()

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}

			for _, st := range s.chats[chatID].userConnections {
				if err := st.Send(msg); err != nil {
					return err
				}
			}

		case <-stream.Context().Done():
			s.chats[chatID].m.Lock()
			delete(s.chats[chatID].userConnections, username)
			s.chats[chatID].m.Unlock()

			return nil
		}
	}
}