package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	serviceMocks "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/mocks"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate(t *testing.T){
	t.Parallel()

	type userServiceMockFunc func(mc *gomock.Controller) service.UserService

	type args struct{
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc = gomock.NewController(t)

		id      = gofakeit.Int64()
		name   = gofakeit.Name()
		username = gofakeit.Username()
		email = gofakeit.Email()
		birthDate = gofakeit.Date()
		avatarUrl = gofakeit.URL()
		

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.UserInfo{
				Name:   name,
				Username: username,
				Email: email,
				BirthDate: timestamppb.New(birthDate),
				AvatarUrl: avatarUrl,
			},
		}

		info = &model.UserInfo{
			    Name:   name,
				Username: username,
				Email: email,
				Birth_date: birthDate,
				Avatar_url: avatarUrl,
		}

		res = &desc.CreateResponse{
			Id: id,
		}

	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *gomock.Controller) service.UserService {
				mock := serviceMocks.NewMockUserService(mc)
				mock.EXPECT().Create(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *gomock.Controller) service.UserService {
				mock := serviceMocks.NewMockUserService(mc)
				mock.EXPECT().Create(ctx, info).Return(int64(0), serviceErr)
				return mock
			},
		},
	}


	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			res, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}

}