package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/mocks"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	mc := gomock.NewController(t)
	ctx := context.Background()

	type mockAction func(mc *gomock.Controller) service.UserService
	type mockAccess func(mc *gomock.Controller) service.AccessService

	correctReq := desc.CreateRequest{
		Info: &desc.UserInfo{
			Role:  desc.UserRole(1),
			Email: "email",
			Name:  "name",
			AvatarUrl: "asasasas",
			Username: "username",
			BirthDate: timestamppb.New(time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC)),
		},
		Pass: &desc.UserPassword{
			Password:        "password",
			PasswordConfirm: "password",
		},
	}

	incorrectReq := desc.CreateRequest{
		Info: &desc.UserInfo{
			Role:  desc.UserRole(1),
			Email: "email",
			Name:  "name",
			AvatarUrl: "asasasas",
			Username: "username",
			BirthDate: timestamppb.New(time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC)),
		},
		Pass: &desc.UserPassword{
			Password:        "password",
			PasswordConfirm: "password12345",
			
			
		},
	}

	paramsCreate := model.UserCreate{
		Info: model.UserInfo{
			Role:  model.UserRole(1),
			Email: "email",
			Name:  "name",
			Avatar_url: "asasasas",
			Username: "username",
			Birth_date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
		},
		Password: "password",
	}

	id := int64(1)

	resp := desc.CreateResponse{
		Id: id,
	}

	tests := []struct {
		name       string
		ctx        context.Context
		req        *desc.CreateRequest
		err        error
		expected   *desc.CreateResponse
		mockAction mockAction
		mockAccess
	}{
		{
			name:     "sucessfull test",
			req:      &correctReq,
			err:      nil,
			ctx:      ctx,
			expected: &resp,
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)
				userServiceMock.EXPECT().
    			Create(ctx, &paramsCreate).
   				Return(id, nil)

				return userServiceMock
			},
			mockAccess: func(mc *gomock.Controller) service.AccessService {
				mock := mocks.NewMockAccessService(mc)

				return mock
			},
		},
		{
			name:     "mismatch passwords",
			req:      &incorrectReq,
			err:      errors.New("passwords mismatch"),
			ctx:      ctx,
			expected: nil,
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)

				return userServiceMock
			},
			mockAccess: func(mc *gomock.Controller) service.AccessService {
				mock := mocks.NewMockAccessService(mc)

				return mock
			},
		},
		{
			name:     "some error",
			req:      &correctReq,
			err:      errors.New("failed to create user: error"),
			ctx:      ctx,
			expected: nil,
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)
				userServiceMock.EXPECT().
    			Create(ctx, &paramsCreate).
   				Return(int64(0), errors.New("error"))

				return userServiceMock
			},
			mockAccess: func(mc *gomock.Controller) service.AccessService {
				mock := mocks.NewMockAccessService(mc)

				return mock
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userServ := test.mockAction(mc)
			accessServ := test.mockAccess(mc)

			impl := user.NewImplementation(userServ, accessServ)

			res, err := impl.Create(test.ctx, test.req)

			require.Equal(t, res, test.expected)
			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}