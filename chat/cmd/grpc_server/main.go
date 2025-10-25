package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/interceptor"
	desc "github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

type server struct{
	desc.UnimplementedChatV1Server
}

func(s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error){

	err:=req.Validate()   //прото валидация
	if err!=nil{
		return nil, err
	}


	log.Printf("Chat id: %d", req.ChatId)
	ids:=[]int64{1,2,3}
	names:=[]string{"first", "second", "third"}
	return &desc.GetResponse{
		ChatId: ids,
		ChatName: names,
		}, nil
}

func main(){
	lis, err :=net.Listen("tcp", fmt.Sprintf(":%d",grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s:=grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),//валидация
	)

	reflection.Register(s)//даем клиекту доступ к ручкам
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}