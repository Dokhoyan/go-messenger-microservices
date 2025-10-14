package accessApi

import (
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
	accesspb "github.com/Dokhoyan/go-messenger-microservices/auth/pkg/api/access_v1"
)

type Implementation struct {
	accesspb.UnimplementedAccessV1Server
	service service.AccessService
}

// NewImplementation - Создает новую имплементацию для gRPC сервера
func NewImplementation(serv service.AccessService) *Implementation {
	return &Implementation{
		service: serv,
	}
}