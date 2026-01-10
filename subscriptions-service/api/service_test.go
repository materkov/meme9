package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
	"go.uber.org/mock/gomock"

	subscriptionsapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/subscriptions"
	"github.com/materkov/meme9/subscriptions-service/api/mocks"
)

func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func reauireError(t *testing.T, err error, msg string) {
	require.Error(t, err)
	twirpErr := err.(twirp.Error)
	require.Equal(t, msg, twirpErr.Msg())
}

func TestService_Subscribe(t *testing.T) {
	ctx := context.Background()

	t.Run("successful subscribe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			Subscribe(gomock.Any(), "user1", "user2").
			Return(nil)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		resp, err := service.Subscribe(ctxWithUser, req)
		require.NoError(t, err)
		require.True(t, resp.Subscribed)
	})

	t.Run("no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		_, err := service.Subscribe(ctx, req)
		reauireError(t, err, "unauthorized")
	})

	t.Run("empty target user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		_, err := service.Subscribe(ctxWithUser, req)
		reauireError(t, err, "target_id_empty")
	})

	t.Run("cannot subscribe to yourself", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user1",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		_, err := service.Subscribe(ctxWithUser, req)
		reauireError(t, err, "cannot_subscribe_to_yourself")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			Subscribe(gomock.Any(), "user1", "user2").
			Return(errors.New("database error"))

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		_, err := service.Subscribe(ctxWithUser, req)
		require.Error(t, err)
		twirpErr := err.(twirp.Error)
		require.Equal(t, twirp.Internal, twirpErr.Code())
		require.Equal(t, "database error", twirpErr.Msg())
	})
}

func TestService_Unsubscribe(t *testing.T) {
	ctx := context.Background()

	t.Run("successful unsubscribe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			Unsubscribe(gomock.Any(), "user1", "user2").
			Return(nil)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		resp, err := service.Unsubscribe(ctxWithUser, req)
		require.NoError(t, err)
		require.False(t, resp.Subscribed)
	})

	t.Run("no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		_, err := service.Unsubscribe(ctx, req)
		reauireError(t, err, "unauthorized")
	})

	t.Run("empty target user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		_, err := service.Unsubscribe(ctxWithUser, req)
		reauireError(t, err, "target_id_empty")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			Unsubscribe(gomock.Any(), "user1", "user2").
			Return(errors.New("database error"))

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		_, err := service.Unsubscribe(ctxWithUser, req)
		require.Error(t, err)
		twirpErr := err.(twirp.Error)
		require.Equal(t, twirp.Internal, twirpErr.Code())
		require.Equal(t, "database error", twirpErr.Msg())
	})
}

func TestService_GetStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get status - subscribed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			IsSubscribed(gomock.Any(), "user1", "user2").
			Return(true, nil)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		resp, err := service.GetStatus(ctxWithUser, req)
		require.NoError(t, err)
		require.True(t, resp.Subscribed)
	})

	t.Run("successful get status - not subscribed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			IsSubscribed(gomock.Any(), "user1", "user2").
			Return(false, nil)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		resp, err := service.GetStatus(ctxWithUser, req)
		require.NoError(t, err)
		require.False(t, resp.Subscribed)
	})

	t.Run("no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		_, err := service.GetStatus(ctx, req)
		reauireError(t, err, "unauthorized")
	})

	t.Run("empty target user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		_, err := service.GetStatus(ctxWithUser, req)
		reauireError(t, err, "target_id_empty")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSubscriptionsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			IsSubscribed(gomock.Any(), "user1", "user2").
			Return(false, errors.New("database error"))

		req := &subscriptionsapi.SubscribeRequest{
			UserId: "user2",
		}

		ctxWithUser := contextWithUserID(ctx, "user1")
		_, err := service.GetStatus(ctxWithUser, req)
		require.Error(t, err)
		twirpErr := err.(twirp.Error)
		require.Equal(t, twirp.Internal, twirpErr.Code())
		require.Equal(t, "database error", twirpErr.Msg())
	})
}
