package api

// To regenerate mocks after changing the MongoAdapter interface, run:
//   go generate ./api/...

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
	"go.uber.org/mock/gomock"

	likesapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/likes"
	"github.com/materkov/meme9/likes-service/api/mocks"
)

func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mocks.NewMockMongoAdapter(ctrl)
	service := NewService(mockAdapter)

	require.NotNil(t, service)
	assert.Equal(t, mockAdapter, service.likes)
}

func TestService_Like(t *testing.T) {
	t.Run("successful like", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			Like(gomock.Any(), "post1", "user1").
			Return(nil)

		service := NewService(mockAdapter)
		ctx := contextWithUserID(context.Background(), "user1")

		req := &likesapi.LikeRequest{PostId: "post1"}
		resp, err := service.Like(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.Liked)
	})

	t.Run("error when user not authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)
		ctx := context.Background() // No user ID in context

		req := &likesapi.LikeRequest{PostId: "post1"}
		resp, err := service.Like(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.Unauthenticated, twirpErr.Code())
		assert.Equal(t, "auth_required", twirpErr.Msg())
	})

	t.Run("error when post_id is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)
		ctx := contextWithUserID(context.Background(), "user1")

		req := &likesapi.LikeRequest{PostId: ""}
		resp, err := service.Like(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.InvalidArgument, twirpErr.Code())
		assert.Equal(t, "post_id_required", twirpErr.Msg())
	})

	t.Run("error when adapter fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			Like(gomock.Any(), "post1", "user1").
			Return(errors.New("database error"))

		service := NewService(mockAdapter)
		ctx := contextWithUserID(context.Background(), "user1")

		req := &likesapi.LikeRequest{PostId: "post1"}
		resp, err := service.Like(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.Internal, twirpErr.Code())
		assert.Contains(t, twirpErr.Msg(), "database error")
	})
}

func TestService_Unlike(t *testing.T) {
	t.Run("successful unlike", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			Unlike(gomock.Any(), "post1", "user1").
			Return(nil)

		service := NewService(mockAdapter)
		ctx := contextWithUserID(context.Background(), "user1")

		req := &likesapi.LikeRequest{PostId: "post1"}
		resp, err := service.Unlike(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.False(t, resp.Liked)
	})

	t.Run("error when user not authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.LikeRequest{PostId: "post1"}
		resp, err := service.Unlike(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.Unauthenticated, twirpErr.Code())
	})

	t.Run("error when post_id is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)
		ctx := contextWithUserID(context.Background(), "user1")

		req := &likesapi.LikeRequest{PostId: ""}
		resp, err := service.Unlike(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.InvalidArgument, twirpErr.Code())
	})

	t.Run("error when adapter fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			Unlike(gomock.Any(), "post1", "user1").
			Return(errors.New("database error"))

		service := NewService(mockAdapter)
		ctx := contextWithUserID(context.Background(), "user1")

		req := &likesapi.LikeRequest{PostId: "post1"}
		resp, err := service.Unlike(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.Internal, twirpErr.Code())
	})
}

func TestService_GetLikers(t *testing.T) {
	t.Run("successful get likers - validates adapter call", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikers(gomock.Any(), "post1", "", 20).
			Return([]string{"user1", "user2"}, "next_token", nil)

		// Note: This test validates the adapter interaction and count logic.
		// The users service integration is tested separately in integration tests.
		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.GetLikersRequest{
			PostId:    "post1",
			PageToken: "",
			Count:     0, // Should default to 20
		}

		resp, err := service.GetLikers(ctx, req)

		// The response will be created even if users service fails (it logs and continues)
		// So we just verify the adapter was called correctly
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "next_token", resp.PageToken)
		// Likers will be populated even if users service fails (with empty username/avatar)
		assert.Len(t, resp.Likers, 2)
	})

	t.Run("error when post_id is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.GetLikersRequest{PostId: ""}
		resp, err := service.GetLikers(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.InvalidArgument, twirpErr.Code())
		assert.Equal(t, "post_id_required", twirpErr.Msg())
	})

	t.Run("count validation - default to 20 when 0", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikers(gomock.Any(), "post1", "", 20).
			Return([]string{}, "", nil)

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.GetLikersRequest{
			PostId: "post1",
			Count:  0,
		}

		_, err := service.GetLikers(ctx, req)
		require.NoError(t, err)
	})

	t.Run("count validation - cap at 100", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikers(gomock.Any(), "post1", "", 100).
			Return([]string{}, "", nil)

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.GetLikersRequest{
			PostId: "post1",
			Count:  150, // Should be capped at 100
		}

		_, err := service.GetLikers(ctx, req)
		require.NoError(t, err)
	})

	t.Run("count validation - use provided count when valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikers(gomock.Any(), "post1", "", 50).
			Return([]string{}, "", nil)

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.GetLikersRequest{
			PostId: "post1",
			Count:  50,
		}

		_, err := service.GetLikers(ctx, req)
		require.NoError(t, err)
	})

	t.Run("error when adapter fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikers(gomock.Any(), "post1", "", 20).
			Return(nil, "", errors.New("database error"))

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.GetLikersRequest{PostId: "post1"}
		resp, err := service.GetLikers(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get likers")
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("handles pagination token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikers(gomock.Any(), "post1", "token123", 20).
			Return([]string{"user3"}, "", nil)

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.GetLikersRequest{
			PostId:    "post1",
			PageToken: "token123",
			Count:     20,
		}

		resp, err := service.GetLikers(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Empty(t, resp.PageToken) // No next page
		assert.Len(t, resp.Likers, 1)
	})
}

func TestService_CheckLike(t *testing.T) {
	t.Run("successful check like", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikedByUser(gomock.Any(), "user1", []string{"post1", "post2"}).
			Return(map[string]bool{
				"post1": true,
				"post2": false,
			}, nil)
		mockAdapter.EXPECT().
			GetLikesCounts(gomock.Any(), []string{"post1", "post2"}).
			Return(map[string]int32{
				"post1": 5,
				"post2": 3,
			}, nil)

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.CheckLikeRequest{
			UserId:  "user1",
			PostIds: []string{"post1", "post2"},
		}

		resp, err := service.CheckLike(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, []bool{true, false}, resp.Liked)
		assert.Equal(t, []int32{5, 3}, resp.LikesCount)
	})

	t.Run("error when user_id is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.CheckLikeRequest{
			UserId:  "",
			PostIds: []string{"post1"},
		}

		resp, err := service.CheckLike(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.InvalidArgument, twirpErr.Code())
		assert.Equal(t, "user_id_required", twirpErr.Msg())
	})

	t.Run("error when post_ids is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.CheckLikeRequest{
			UserId:  "user1",
			PostIds: []string{},
		}

		resp, err := service.CheckLike(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		assert.Equal(t, twirp.InvalidArgument, twirpErr.Code())
		assert.Equal(t, "post_ids_required", twirpErr.Msg())
	})

	t.Run("error when getLikedByUser fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikedByUser(gomock.Any(), "user1", []string{"post1"}).
			Return(nil, errors.New("database error"))

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.CheckLikeRequest{
			UserId:  "user1",
			PostIds: []string{"post1"},
		}

		resp, err := service.CheckLike(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get liked status")
	})

	t.Run("error when getLikesCounts fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikedByUser(gomock.Any(), "user1", []string{"post1"}).
			Return(map[string]bool{"post1": true}, nil)
		mockAdapter.EXPECT().
			GetLikesCounts(gomock.Any(), []string{"post1"}).
			Return(nil, errors.New("database error"))

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.CheckLikeRequest{
			UserId:  "user1",
			PostIds: []string{"post1"},
		}

		resp, err := service.CheckLike(ctx, req)

		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get likes counts")
	})

	t.Run("preserves order of postIds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockMongoAdapter(ctrl)
		mockAdapter.EXPECT().
			GetLikedByUser(gomock.Any(), "user1", []string{"post1", "post2", "post3"}).
			Return(map[string]bool{
				"post3": true,
				"post1": false,
				"post2": true,
			}, nil)
		mockAdapter.EXPECT().
			GetLikesCounts(gomock.Any(), []string{"post1", "post2", "post3"}).
			Return(map[string]int32{
				"post1": 1,
				"post2": 2,
				"post3": 3,
			}, nil)

		service := NewService(mockAdapter)
		ctx := context.Background()

		req := &likesapi.CheckLikeRequest{
			UserId:  "user1",
			PostIds: []string{"post1", "post2", "post3"},
		}

		resp, err := service.CheckLike(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		// Should preserve order: post1, post2, post3
		assert.Equal(t, []bool{false, true, true}, resp.Liked)
		assert.Equal(t, []int32{1, 2, 3}, resp.LikesCount)
	})
}
