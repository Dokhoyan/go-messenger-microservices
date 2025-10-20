package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/mocks"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGet(t *testing.T) {
	t.Parallel()

	mc := gomock.NewController(t)
	ctx := context.Background()

	type mockAction func(mc *gomock.Controller) service.UserService
	type mockAccess func(mc *gomock.Controller) service.AccessService

	correctReq := &desc.GetRequest{
		Id: 1,
	}

	id := int64(1)

	timeNow := time.Now()

	resGet := &model.User{
		ID: id,
		Info: model.UserInfo{
			Role:  model.UserRole(1),
			Email: "email",
			Name:  "name",
			Avatar_url: "asasasas",
			Username: "username",
			Birth_date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
		},
		UpdatedAt: sql.NullTime{
			Time:  timeNow,
			Valid: true,
		},
		CreatedAt: timeNow,
	}

	resp := &desc.GetResponse{
		User: &desc.User{
			Id: id,
			Info: &desc.UserInfo{
				Role:  desc.UserRole(1),
				Email: "email",
				Name:  "name",
				AvatarUrl: "asasasas",
				Username: "username",
				BirthDate: timestamppb.New(time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC)),
			},
			CreatedAt: timestamppb.New(timeNow),
			UpdatedAt: timestamppb.New(timeNow),
		},
	}

	tests := []struct {
		name       string
		ctx        context.Context
		req        *desc.GetRequest
		err        error
		expected   *desc.GetResponse
		mockAction mockAction
		mockAccess
	}{
		{
			name:     "sucessfull test",
			req:      correctReq,
			err:      nil,
			ctx:      ctx,
			expected: resp,
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)
				userServiceMock.EXPECT().
				Get(ctx, id).
				Return(resGet, nil)
				return userServiceMock
			},
			mockAccess: func(mc *gomock.Controller) service.AccessService {
				mock := mocks.NewMockAccessService(mc)

				return mock
			},
		},
		{
			name:     "some error",
			req:      correctReq,
			err:      errors.New("failed to get user: error"),
			ctx:      ctx,
			expected: nil,
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)
				userServiceMock.EXPECT().
				Get(ctx, id).
				Return(nil, errors.New("error"))
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
			

			mockAuth := mocks.NewMockAuthDataService(mc)
			impl := user.NewImplementation(userServ, mockAuth)

			res, err := impl.Get(test.ctx, test.req)

			require.Equal(t, res, test.expected)
			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}