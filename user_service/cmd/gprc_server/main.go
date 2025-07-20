package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

type server struct{
	desc.UnimplementedUserV1Server
}

func(s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error){

	err:=req.Validate()   //прото валидация
	if err!=nil{
		return nil, err
	}


	return nil, nil
}

func main(){
	lis, err :=net.Listen("tcp", fmt.Sprintf(":%d",grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s:=grpc.NewServer()//валидация иеще что то 

	reflection.Register(s)//даем клиекту доступ к ручкам
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}