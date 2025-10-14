package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/api/user"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/service/mocks"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	mc := gomock.NewController(t)
	ctx := context.Background()

	type mockAction func(mc *gomock.Controller) service.UserService
	type mockAccess func(mc *gomock.Controller) service.AccessService

	id := int64(1)

	correctReq := &desc.DeleteRequest{
		Id: id,
	}

	tests := []struct {
		name       string
		ctx        context.Context
		req        *desc.DeleteRequest
		err        error
		expected   *empty.Empty
		mockAction mockAction
		mockAccess
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
				Delete(ctx, id).
				Return(nil)
				return userServiceMock
			},
			mockAccess: func(mc *gomock.Controller) service.AccessService {
				mock := mocks.NewMockAccessService(mc)
				mock.EXPECT().Check(gomock.Any(), gomock.Any()).Return(nil)

				return mock
			},
		},
		{
			name:     "some error",
			req:      correctReq,
			err:      errors.New("failed to delete user: error"),
			ctx:      ctx,
			expected: nil,
			mockAction: func(mc *gomock.Controller) service.UserService {
				userServiceMock := mocks.NewMockUserService(mc)
				userServiceMock.EXPECT().
				Delete(ctx, id).
				Return(errors.New("error"))

				return userServiceMock
			},
			mockAccess: func(mc *gomock.Controller) service.AccessService {
				mock := mocks.NewMockAccessService(mc)
				mock.EXPECT().Check(gomock.Any(), gomock.Any()).Return(nil)

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

			mockAuth := mocks.NewMockAuthDataService(mc)
			impl := user.NewImplementation(userServ, accessServ, mockAuth)

			res, err := impl.Delete(test.ctx, test.req)

			require.Equal(t, res, test.expected)
			if err != nil && test.err != nil {
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}