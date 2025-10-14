package authApi

import (
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/service"
	desc "github.com/Dokhoyan/go-messenger-microservices/auth/pkg/api/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation - создает новую имплементацию для gRPC сервера
func NewImplementation(service service.AuthService) *Implementation {
	return &Implementation{
		authService: service,
	}
}