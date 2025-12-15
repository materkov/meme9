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

func initService(ctrl *gomock.Controller) (*Service, *mocks.MockUsersAdapter) {
	mockUsers := mocks.NewMockUsersAdapter(ctrl)
	return NewService(mockUsers), mockUsers
}

func TestService_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockUsers := initService(ctrl)

	ctx := context.Background()
	userID := "user-123"
	username := "testuser"

	user := &users.User{
		ID:       userID,
		Username: username,
	}

	mockUsers.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil).
		Times(1)

	resp, err := service.Get(ctx, &usersapi.GetUserRequest{
		UserId: userID,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, userID, resp.Id)
	require.Equal(t, username, resp.Username)
}

func TestService_Get_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, _ := initService(ctrl)
	ctx := context.Background()

	t.Run("empty user id", func(t *testing.T) {
		_, err := service.Get(ctx, &usersapi.GetUserRequest{
			UserId: "",
		})
		api.RequireError(t, err, "user_not_found")
	})
}

func TestService_Get_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockUsers := initService(ctrl)
	ctx := context.Background()
	userID := "user-not-found"

	mockUsers.EXPECT().
		GetByID(ctx, userID).
		Return(nil, users.ErrNotFound).
		Times(1)

	resp, err := service.Get(ctx, &usersapi.GetUserRequest{
		UserId: userID,
	})

	require.Error(t, err)
	require.Nil(t, resp)
	api.RequireError(t, err, "user_not_found")
}
