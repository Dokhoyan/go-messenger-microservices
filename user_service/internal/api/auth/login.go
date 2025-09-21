package auth

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/auth/converter"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/auth_v1"
	"github.com/pkg/errors"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, converter.AuthProtoToAuthDTO(req))
	if err != nil {
		return nil, errors.Errorf("authentification error: %s", err)
	}


	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, err
}