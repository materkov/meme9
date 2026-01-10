package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
	"go.uber.org/mock/gomock"

	postsapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/posts"
	"github.com/materkov/meme9/posts-service/adapters/posts"
	"github.com/materkov/meme9/posts-service/api/mocks"
)

func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func reauireError(t *testing.T, err error, msg string) {
	require.Error(t, err)
	twirpErr := err.(twirp.Error)
	require.Equal(t, msg, twirpErr.Msg())
}

func TestService_Publish(t *testing.T) {
	ctx := context.Background()

	t.Run("successful publish", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		post := &posts.Post{
			ID:        "post123",
			Text:      "Test post",
			UserID:    "user123",
			CreatedAt: time.Now(),
		}

		mockAdapter.EXPECT().
			Add(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, p posts.Post) (*posts.Post, error) {
				require.Equal(t, "Test post", p.Text)
				require.Equal(t, "user123", p.UserID)
				require.NotZero(t, p.CreatedAt)
				return post, nil
			})

		req := &postsapi.PublishRequest{
			Text: "Test post",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		resp, err := service.Publish(ctxWithUser, req)
		require.NoError(t, err)
		require.Equal(t, "post123", resp.Id)
	})

	t.Run("no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &postsapi.PublishRequest{
			Text: "Test post",
		}

		_, err := service.Publish(ctx, req)
		reauireError(t, err, "auth_required")
	})

	t.Run("empty text", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &postsapi.PublishRequest{
			Text: "",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Publish(ctxWithUser, req)
		reauireError(t, err, "text_empty")
	})

	t.Run("text too long", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		longText := make([]byte, 1001)
		for i := range longText {
			longText[i] = 'a'
		}

		req := &postsapi.PublishRequest{
			Text: string(longText),
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Publish(ctxWithUser, req)
		reauireError(t, err, "text_too_long")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			Add(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("database error"))

		req := &postsapi.PublishRequest{
			Text: "Test post",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Publish(ctxWithUser, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create post")
	})
}

func TestService_GetByUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get by users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		postsList := []posts.Post{
			{
				ID:        "post1",
				Text:      "Post 1",
				UserID:    "user1",
				CreatedAt: time.Now(),
			},
			{
				ID:        "post2",
				Text:      "Post 2",
				UserID:    "user1",
				CreatedAt: time.Now(),
			},
		}

		mockAdapter.EXPECT().
			GetByUserIDs(gomock.Any(), []string{"user1"}).
			Return(postsList, nil)

		req := &postsapi.GetByUsersRequest{
			UserId: "user1",
		}

		resp, err := service.GetByUsers(ctx, req)
		require.NoError(t, err)
		require.Len(t, resp.Posts, 2)
		require.Equal(t, "post1", resp.Posts[0].Id)
		require.Equal(t, "Post 1", resp.Posts[0].Text)
	})

	t.Run("empty user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &postsapi.GetByUsersRequest{
			UserId: "",
		}

		_, err := service.GetByUsers(ctx, req)
		reauireError(t, err, "user_id_required")
	})

	t.Run("adapter error returns empty list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByUserIDs(gomock.Any(), []string{"user1"}).
			Return(nil, errors.New("database error"))

		req := &postsapi.GetByUsersRequest{
			UserId: "user1",
		}

		resp, err := service.GetByUsers(ctx, req)
		require.NoError(t, err)
		require.Empty(t, resp.Posts)
	})

	t.Run("empty posts list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByUserIDs(gomock.Any(), []string{"user1"}).
			Return([]posts.Post{}, nil)

		req := &postsapi.GetByUsersRequest{
			UserId: "user1",
		}

		resp, err := service.GetByUsers(ctx, req)
		require.NoError(t, err)
		require.Empty(t, resp.Posts)
	})
}

func TestService_Get(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		post := &posts.Post{
			ID:        "post123",
			Text:      "Test post",
			UserID:    "user123",
			CreatedAt: time.Now(),
			Deleted:   false,
		}

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(post, nil)

		req := &postsapi.GetPostRequest{
			PostId: "post123",
		}

		resp, err := service.Get(ctx, req)
		require.NoError(t, err)
		require.Equal(t, "post123", resp.Id)
		require.Equal(t, "Test post", resp.Text)
		require.Equal(t, "user123", resp.UserId)
	})

	t.Run("empty post ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &postsapi.GetPostRequest{
			PostId: "",
		}

		_, err := service.Get(ctx, req)
		reauireError(t, err, "post_id_required")
	})

	t.Run("post not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(nil, posts.ErrNotFound)

		req := &postsapi.GetPostRequest{
			PostId: "post123",
		}

		_, err := service.Get(ctx, req)
		reauireError(t, err, "post_not_found")
	})

	t.Run("deleted post", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		post := &posts.Post{
			ID:        "post123",
			Text:      "Deleted post",
			UserID:    "user123",
			CreatedAt: time.Now(),
			Deleted:   true,
		}

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(post, nil)

		req := &postsapi.GetPostRequest{
			PostId: "post123",
		}

		_, err := service.Get(ctx, req)
		reauireError(t, err, "post_not_found")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(nil, errors.New("database error"))

		req := &postsapi.GetPostRequest{
			PostId: "post123",
		}

		_, err := service.Get(ctx, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get post")
	})
}

func TestService_GetFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get feed all", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		postsList := []posts.Post{
			{
				ID:        "post1",
				Text:      "Post 1",
				UserID:    "user1",
				CreatedAt: time.Now(),
			},
		}

		mockAdapter.EXPECT().
			GetAll(gomock.Any()).
			Return(postsList, nil)

		req := &postsapi.FeedRequest{
			Type: postsapi.FeedType_FEED_TYPE_ALL,
		}

		resp, err := service.GetFeed(ctx, req)
		require.NoError(t, err)
		require.Len(t, resp.Posts, 1)
		require.Equal(t, "post1", resp.Posts[0].Id)
	})

	t.Run("unspecified feed type defaults to all", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		postsList := []posts.Post{}

		mockAdapter.EXPECT().
			GetAll(gomock.Any()).
			Return(postsList, nil)

		req := &postsapi.FeedRequest{
			Type: postsapi.FeedType_FEED_TYPE_UNSPECIFIED,
		}

		_, err := service.GetFeed(ctx, req)
		require.NoError(t, err)
	})

	t.Run("subscriptions feed requires auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &postsapi.FeedRequest{
			Type: postsapi.FeedType_FEED_TYPE_SUBSCRIPTIONS,
		}

		_, err := service.GetFeed(ctx, req)
		reauireError(t, err, "auth_required")
	})

	t.Run("adapter error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetAll(gomock.Any()).
			Return(nil, errors.New("database error"))

		req := &postsapi.FeedRequest{
			Type: postsapi.FeedType_FEED_TYPE_ALL,
		}

		_, err := service.GetFeed(ctx, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get posts")
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		post := &posts.Post{
			ID:        "post123",
			Text:      "Test post",
			UserID:    "user123",
			CreatedAt: time.Now(),
		}

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(post, nil)

		mockAdapter.EXPECT().
			MarkAsDeleted(gomock.Any(), "post123").
			Return(nil)

		req := &postsapi.DeleteRequest{
			PostId: "post123",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Delete(ctxWithUser, req)
		require.NoError(t, err)
	})

	t.Run("no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &postsapi.DeleteRequest{
			PostId: "post123",
		}

		_, err := service.Delete(ctx, req)
		reauireError(t, err, "auth_required")
	})

	t.Run("empty post ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		req := &postsapi.DeleteRequest{
			PostId: "",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Delete(ctxWithUser, req)
		reauireError(t, err, "post_id_required")
	})

	t.Run("post not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(nil, posts.ErrNotFound)

		req := &postsapi.DeleteRequest{
			PostId: "post123",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Delete(ctxWithUser, req)
		reauireError(t, err, "post_not_found")
	})

	t.Run("not post owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		post := &posts.Post{
			ID:        "post123",
			Text:      "Test post",
			UserID:    "user456",
			CreatedAt: time.Now(),
		}

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(post, nil)

		req := &postsapi.DeleteRequest{
			PostId: "post123",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Delete(ctxWithUser, req)
		reauireError(t, err, "not_post_owner")
	})

	t.Run("error getting post", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(nil, errors.New("database error"))

		req := &postsapi.DeleteRequest{
			PostId: "post123",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Delete(ctxWithUser, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get post")
	})

	t.Run("error marking as deleted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockPostsAdapter(ctrl)
		service := NewService(mockAdapter)

		post := &posts.Post{
			ID:        "post123",
			Text:      "Test post",
			UserID:    "user123",
			CreatedAt: time.Now(),
		}

		mockAdapter.EXPECT().
			GetByID(gomock.Any(), "post123").
			Return(post, nil)

		mockAdapter.EXPECT().
			MarkAsDeleted(gomock.Any(), "post123").
			Return(errors.New("delete error"))

		req := &postsapi.DeleteRequest{
			PostId: "post123",
		}

		ctxWithUser := contextWithUserID(ctx, "user123")
		_, err := service.Delete(ctxWithUser, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to delete post")
	})
}
