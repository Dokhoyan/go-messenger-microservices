package user

import (
	"context"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
)

func (i *Implementation) GetUserAuthData(ctx context.Context, req *user_v1.GetUserAuthDataRequest) (*user_v1.GetUserAuthDataResponse, error) {
    user, err := i.authDataService.GetUserByUsername(ctx, req.GetUsername())
    if err != nil {
        return nil, err
    }

    return &user_v1.GetUserAuthDataResponse{
        User: &user_v1.UserAuthData{
            Id:           user.ID,
            Username:     user.Info.Username,
            PasswordHash: user.Password,
			Role: user_v1.UserRole(user.Info.Role),
        },
    }, nil
}