package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
	"go.uber.org/mock/gomock"

	usersapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/users"
	"github.com/materkov/meme9/users-service/adapters/users"
	"github.com/materkov/meme9/users-service/api/mocks"
)

func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func reauireError(t *testing.T, err error, msg string) {
	require.Error(t, err)
	twirpErr := err.(twirp.Error)
	require.Equal(t, msg, twirpErr.Msg())
}

func TestService_Get(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		user := &users.User{
			ID:        "user123",
			Username:  "testuser",
			AvatarURL: "https://example.com/avatar.jpg",
		}

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "user123").
			Return(user, nil)

		req := &usersapi.GetUserRequest{
			UserId: "user123",
		}

		resp, err := service.Get(ctx, req)
		require.NoError(t, err)
		require.Equal(t, "user123", resp.Id)
		require.Equal(t, "testuser", resp.Username)
		require.Equal(t, "https://example.com/avatar.jpg", resp.AvatarUrl)
	})

	t.Run("empty user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &usersapi.GetUserRequest{
			UserId: "",
		}

		_, err := service.Get(ctx, req)
		reauireError(t, err, "user_not_found")
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "user123").
			Return(nil, users.ErrNotFound)

		req := &usersapi.GetUserRequest{
			UserId: "user123",
		}

		_, err := service.Get(ctx, req)
		reauireError(t, err, "user_not_found")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "user123").
			Return(nil, errors.New("database error"))

		req := &usersapi.GetUserRequest{
			UserId: "user123",
		}

		_, err := service.Get(ctx, req)
		reauireError(t, err, "user_not_found")
	})
}

func TestService_SetAvatar(t *testing.T) {
	ctx := context.Background()

	t.Run("successful set avatar", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			UpdateAvatar(gomock.Any(), "user123", "https://example.com/avatar.jpg").
			Return(nil)

		req := &usersapi.SetAvatarRequest{
			UserId:    "user123",
			AvatarUrl: "https://example.com/avatar.jpg",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.SetAvatar(ctxWithUser, req)
		require.NoError(t, err)
	})

	t.Run("no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &usersapi.SetAvatarRequest{
			UserId:    "user123",
			AvatarUrl: "https://example.com/avatar.jpg",
		}

		_, err := service.SetAvatar(ctx, req)
		reauireError(t, err, "auth_required")
	})

	t.Run("empty user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &usersapi.SetAvatarRequest{
			UserId:    "",
			AvatarUrl: "https://example.com/avatar.jpg",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.SetAvatar(ctxWithUser, req)
		reauireError(t, err, "user_id_required")
	})

	t.Run("cannot set other user's avatar", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &usersapi.SetAvatarRequest{
			UserId:    "user456",
			AvatarUrl: "https://example.com/avatar.jpg",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.SetAvatar(ctxWithUser, req)
		reauireError(t, err, "can_only_set_own_avatar")
	})

	t.Run("empty avatar URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &usersapi.SetAvatarRequest{
			UserId:    "user123",
			AvatarUrl: "",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.SetAvatar(ctxWithUser, req)
		reauireError(t, err, "avatar_url_required")
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			UpdateAvatar(gomock.Any(), "user123", "https://example.com/avatar.jpg").
			Return(users.ErrNotFound)

		req := &usersapi.SetAvatarRequest{
			UserId:    "user123",
			AvatarUrl: "https://example.com/avatar.jpg",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.SetAvatar(ctxWithUser, req)
		reauireError(t, err, "user_not_found")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockUsersAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			UpdateAvatar(gomock.Any(), "user123", "https://example.com/avatar.jpg").
			Return(errors.New("database error"))

		req := &usersapi.SetAvatarRequest{
			UserId:    "user123",
			AvatarUrl: "https://example.com/avatar.jpg",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.SetAvatar(ctxWithUser, req)
		require.Error(t, err)
		twirpErr := err.(twirp.Error)
		require.Equal(t, twirp.Internal, twirpErr.Code())
		require.Contains(t, twirpErr.Msg(), "failed to update avatar")
	})
}
