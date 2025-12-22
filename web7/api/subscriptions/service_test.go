package subscriptions

import (
	"context"
	"testing"

	"github.com/materkov/meme9/web7/api"
	"github.com/materkov/meme9/web7/api/subscriptions/mocks"
	proto "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/subscriptions"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func initService(t *testing.T) (*Service, *mocks.MockSubscriptionsAdapter, func()) {
	ctrl := gomock.NewController(t)
	closer := func() {
		ctrl.Finish()
	}

	mockSubscriptions := mocks.NewMockSubscriptionsAdapter(ctrl)
	return NewService(mockSubscriptions), mockSubscriptions, closer
}

func TestService_Subscribe(t *testing.T) {
	service, mockSubscriptions, closer := initService(t)
	defer closer()

	t.Run("success", func(t *testing.T) {
		mockSubscriptions.EXPECT().
			Subscribe(gomock.Any(), "user-1", "user-2").
			Return(nil)

		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-1")
		resp, err := service.Subscribe(ctx, &proto.SubscribeRequest{
			UserId: "user-2",
		})
		require.NoError(t, err)
		require.True(t, resp.Subscribed)
	})

	t.Run("empty user id", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.Subscribe(ctx, &proto.SubscribeRequest{
			UserId: "user-2",
		})
		api.RequireError(t, err, "unauthorized")
	})

	t.Run("empty target id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-1")
		_, err := service.Subscribe(ctx, &proto.SubscribeRequest{
			UserId: "",
		})
		api.RequireError(t, err, "target_id_empty")
	})

	t.Run("sub to myself", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-1")
		_, err := service.Subscribe(ctx, &proto.SubscribeRequest{
			UserId: "user-1",
		})
		api.RequireError(t, err, "cannot_subscribe_to_yourself")
	})
}

func TestService_Unsubscribe(t *testing.T) {
	service, mockSubscriptions, closer := initService(t)
	defer closer()

	t.Run("success", func(t *testing.T) {
		mockSubscriptions.EXPECT().
			Unsubscribe(gomock.Any(), "user-1", "user-2").
			Return(nil)

		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-1")
		resp, err := service.Unsubscribe(ctx, &proto.SubscribeRequest{
			UserId: "user-2",
		})
		require.NoError(t, err)
		require.False(t, resp.Subscribed)
	})

	t.Run("empty user id", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.Unsubscribe(ctx, &proto.SubscribeRequest{
			UserId: "user-2",
		})
		api.RequireError(t, err, "unauthorized")
	})

	t.Run("empty target id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-1")
		_, err := service.Unsubscribe(ctx, &proto.SubscribeRequest{
			UserId: "",
		})
		api.RequireError(t, err, "target_id_empty")
	})
}

func TestService_GetStatus(t *testing.T) {
	service, mockSubscriptions, closer := initService(t)
	defer closer()

	t.Run("success", func(t *testing.T) {
		mockSubscriptions.EXPECT().
			IsSubscribed(gomock.Any(), "user-1", "user-2").
			Return(true, nil)

		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-1")
		resp, err := service.GetStatus(ctx, &proto.SubscribeRequest{
			UserId: "user-2",
		})
		require.NoError(t, err)
		require.True(t, resp.Subscribed)
	})

	t.Run("empty user id", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetStatus(ctx, &proto.SubscribeRequest{
			UserId: "user-2",
		})
		api.RequireError(t, err, "unauthorized")
	})

	t.Run("empty target id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), api.UserIDKey, "user-1")
		_, err := service.GetStatus(ctx, &proto.SubscribeRequest{
			UserId: "",
		})
		api.RequireError(t, err, "target_id_empty")
	})
}
