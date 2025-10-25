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
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	mc := gomock.NewController(t)
	ctx := context.Background()

	type mockAction func(mc *gomock.Controller) service.UserService

	id := int64(1)

	correctReq := &desc.UpdateRequest{
		Id: id,
		Info: &desc.UpdateInfo{
			Name: &wrappers.StringValue{
				Value: "name",
			},
			BirthDate: timestamppb.New(time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC)),

			AvatarUrl: &wrappers.StringValue{
				Value: "avatar_url",
			},
			Username: &wrappers.StringValue{
				Value: "username",
			},
			Email: &wrappers.StringValue{
				Value: "email",
			},
			Role: desc.UserRole(1),
		},
	}

	updateDTOData := &model.UserUpdate{
		Info: model.UserInfo{
			Role:  model.UserRole(1),
			Email: "email",
			Name:  "name",
			Avatar_url: "avatar_url",
			Birth_date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			Username: "username",
		},
		ID: id,
	}

	tests := []struct {
		name       string
		ctx        context.Context
		req        *desc.UpdateRequest
		err        error
		expected   *empty.Empty
		mockAction mockAction
	}{
		{
			name:     "sucessfull test",
			req:      correctReq,
			err:      nil,
			ctx:      ctx,
			expected: &empty.Empty{},
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)
				userServiceMock.EXPECT().
				Update(ctx, updateDTOData).
				Return(nil)
				return userServiceMock
			},
		},
		{
			name:     "some error",
			req:      correctReq,
			err:      errors.New("failed to update user: error"),
			ctx:      ctx,
			expected: nil,
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)
				userServiceMock.EXPECT().
				Update(ctx, updateDTOData).
				Return(errors.New("error"))
				return userServiceMock
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userServ := test.mockAction(mc)

			impl := user.NewImplementation(userServ)

			res, err := impl.Update(test.ctx, test.req)

			require.Equal(t, res, test.expected)
			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}