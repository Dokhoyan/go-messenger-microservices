package authApi

import (
	"context"

	desc "github.com/Dokhoyan/go-messenger-microservices/auth/pkg/api/auth_v1"
	"github.com/pkg/errors"
)

// GetRefreshToken - возвращает новый refresh токен
func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, errors.Errorf("refresh token update error: %s", err)
	}

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}

// GetAccessToken - возвращает новый access токен
func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, errors.Errorf("access token get error: %s", err)
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}

