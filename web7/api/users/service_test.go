package users

import (
	"context"
	"testing"

	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	"github.com/materkov/meme9/web7/api/users/mocks"
	usersapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/users"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func initService(t *testing.T) (*Service, *mocks.MockUsersAdapter, func()) {
	ctrl := gomock.NewController(t)
	closer := func() {
		ctrl.Finish()
	}
	defer ctrl.Finish()

	mockUsers := mocks.NewMockUsersAdapter(ctrl)
	return NewService(mockUsers), mockUsers, closer
}

func TestService_Get(t *testing.T) {
	service, mockUsers, closer := initService(t)
	defer closer()

	t.Run("success", func(t *testing.T) {
		mockUsers.EXPECT().
			GetByID(gomock.Any(), "user-123").
			Return(&users.User{
				ID:       "user-123",
				Username: "testuser",
			}, nil).
			Times(1)

		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-123")
		resp, err := service.Get(ctx, &usersapi.GetUserRequest{
			UserId: "user-123",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "user-123", resp.Id)
		require.Equal(t, "testuser", resp.Username)
	})

	t.Run("empty user id", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.Get(ctx, &usersapi.GetUserRequest{
			UserId: "",
		})
		api.RequireError(t, err, "user_not_found")
	})

	t.Run("user not found", func(t *testing.T) {
		mockUsers.EXPECT().
			GetByID(gomock.Any(), "user-123").
			Return(nil, users.ErrNotFound).
			Times(1)
		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-123")
		_, err := service.Get(ctx, &usersapi.GetUserRequest{
			UserId: "user-123",
		})
		api.RequireError(t, err, "user_not_found")
	})
}
