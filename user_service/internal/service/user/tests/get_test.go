package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	repoMocks "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/mocks"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/user"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *gomock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = gomock.NewController(t)

		id      = gofakeit.Int64()
		name   = gofakeit.Name()
		username = gofakeit.Username()
		email = gofakeit.Email()
		birthDate = gofakeit.Date()
		avatarUrl = gofakeit.URL()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")

		res = &model.User{
			ID: id,
			Info: model.UserInfo{
				Name:   name,
				Username: username,
				Email: email,
				Birth_date: birthDate,
				Avatar_url: avatarUrl,
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *gomock.Controller) repository.UserRepository {
				mock := repoMocks.NewMockUserRepository(mc)
				mock.EXPECT().Get(ctx, id).Return(res, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  repoErr,
			userRepositoryMock: func(mc *gomock.Controller) repository.UserRepository {
				mock := repoMocks.NewMockUserRepository(mc)
				mock.EXPECT().Get(ctx, id).Return(nil,repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			service := user.NewMockService(userRepoMock)

			newID, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}