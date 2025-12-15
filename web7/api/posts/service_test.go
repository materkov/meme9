package posts

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	"github.com/materkov/meme9/web7/api/posts/mocks"
	postsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/posts"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func contextWithUserID(userID string) context.Context {
	return context.WithValue(context.Background(), api.UserIDKey, userID)
}

func initService(ctrl *gomock.Controller) (*Service, *mocks.MockPostsAdapter, *mocks.MockUsersAdapter, *mocks.MockSubscriptionsAdapter) {
	mockPosts := mocks.NewMockPostsAdapter(ctrl)
	mockUsers := mocks.NewMockUsersAdapter(ctrl)
	mockSubscriptions := mocks.NewMockSubscriptionsAdapter(ctrl)

	return NewService(mockPosts, mockUsers, mockSubscriptions), mockPosts, mockUsers, mockSubscriptions
}

func TestService_Publish_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockPosts, mockUsers, _ := initService(ctrl)
	mockPosts.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(&posts.Post{ID: "post-123"}, nil).
		Times(1)
	mockPosts.EXPECT().
		GetByID(gomock.Any(), "post-123").
		Return(&posts.Post{ID: "post-123"}, nil).
		Times(1)
	mockUsers.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		Return(&users.User{ID: "user-123"}, nil).
		Times(1)

	ctx := contextWithUserID("user-123")
	respPublish, err := service.Publish(ctx, &postsapi.PublishRequest{
		Text: "Test post text",
	})
	require.NoError(t, err)
	require.NotEmpty(t, respPublish.Id)

	respGet, err := service.Get(ctx, &postsapi.GetPostRequest{
		PostId: respPublish.Id,
	})
	require.NoError(t, err)
	require.Equal(t, respPublish.Id, respGet.Id)
	require.NotEmpty(t, respGet.CreatedAt)
}

func TestService_Publish_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := contextWithUserID("user-123")
	service, _, _, _ := initService(ctrl)

	t.Run("no auth", func(t *testing.T) {
		_, err := service.Publish(context.Background(), &postsapi.PublishRequest{
			Text: "Test post text",
		})
		api.RequireError(t, err, "auth_required")
	})

	t.Run("empty text", func(t *testing.T) {
		_, err := service.Publish(ctx, &postsapi.PublishRequest{
			Text: "",
		})
		api.RequireError(t, err, "text_empty")
	})

	t.Run("text too long", func(t *testing.T) {
		_, err := service.Publish(ctx, &postsapi.PublishRequest{
			Text: strings.Repeat("a", 1001),
		})
		api.RequireError(t, err, "text_too_long")
	})
}

func TestService_GetByUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockPosts, mockUsers, _ := initService(ctrl)

	ctx := context.Background()
	userID := "user-123"
	username := "testuser"

	postsList := []posts.Post{
		{ID: "post-1", Text: "Post 1", UserID: userID, CreatedAt: time.Now()},
		{ID: "post-2", Text: "Post 2", UserID: userID, CreatedAt: time.Now()},
	}
	mockPosts.EXPECT().
		GetByUserID(ctx, userID).
		Return(postsList, nil).Times(1)

	user := &users.User{
		ID:       userID,
		Username: username,
	}
	mockUsers.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil).Times(1)

	resp, err := service.GetByUsers(ctx, &postsapi.GetByUsersRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	require.Len(t, resp.Posts, 2)
	require.Equal(t, "post-1", resp.Posts[0].Id)
	require.Equal(t, "post-2", resp.Posts[1].Id)
}

func TestService_GetByUsers_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, _, _, _ := initService(ctrl)

	ctx := context.Background()

	t.Run("empty user id", func(t *testing.T) {
		_, err := service.GetByUsers(ctx, &postsapi.GetByUsersRequest{
			UserId: "",
		})
		api.RequireError(t, err, "user_id_required")
	})
}

func TestService_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockPosts, mockUsers, _ := initService(ctrl)

	ctx := context.Background()
	postID := "post-123"
	userID := "user-123"
	username := "testuser"

	post := &posts.Post{
		ID:        postID,
		Text:      "Test post",
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	mockPosts.EXPECT().
		GetByID(ctx, postID).
		Return(post, nil)

	user := &users.User{
		ID:       userID,
		Username: username,
	}
	mockUsers.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil)

	resp, err := service.Get(ctx, &postsapi.GetPostRequest{
		PostId: postID,
	})

	require.NoError(t, err)
	require.Equal(t, postID, resp.Id)
}

func TestService_Get_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockPosts, _, _ := initService(ctrl)
	mockPosts.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		Return(nil, posts.ErrNotFound).
		Times(1)
	ctx := context.Background()

	t.Run("empty post id", func(t *testing.T) {
		_, err := service.Get(ctx, &postsapi.GetPostRequest{
			PostId: "",
		})
		api.RequireError(t, err, "post_id_required")
	})

	t.Run("post not found", func(t *testing.T) {
		_, err := service.Get(ctx, &postsapi.GetPostRequest{
			PostId: "post-not-found",
		})
		api.RequireError(t, err, "post_not_found")
	})
}

func TestService_GetFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockPosts, mockUsers, _ := initService(ctrl)

	ctx := context.Background()
	postsList := []posts.Post{
		{ID: "post-1", Text: "Post 1", UserID: "user-1", CreatedAt: time.Now()},
		{ID: "post-2", Text: "Post 2", UserID: "user-2", CreatedAt: time.Now()},
	}
	mockPosts.EXPECT().
		GetAll(ctx).
		Return(postsList, nil).Times(1)

	usersMap := map[string]*users.User{
		"user-1": {ID: "user-1", Username: "user1"},
		"user-2": {ID: "user-2", Username: "user2"},
	}
	mockUsers.EXPECT().
		GetByIDs(ctx, []string{"user-1", "user-2"}).
		Return(usersMap, nil).Times(1)

	resp, err := service.GetFeed(ctx, &postsapi.FeedRequest{
		Type: postsapi.FeedType_FEED_TYPE_ALL,
	})

	require.NoError(t, err)
	require.Len(t, resp.Posts, 2)
	require.Equal(t, "post-1", resp.Posts[0].Id)
	require.Equal(t, "post-2", resp.Posts[1].Id)
}

func TestService_GetFeed_UnspecifiedDefaultsToAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockPosts, mockUsers, _ := initService(ctrl)

	ctx := context.Background()
	postsList := []posts.Post{
		{ID: "post-1", Text: "Post 1", UserID: "user-1", CreatedAt: time.Now()},
	}
	mockPosts.EXPECT().
		GetAll(ctx).
		Return(postsList, nil)

	mockUsers.EXPECT().
		GetByIDs(ctx, []string{"user-1"}).
		Return(map[string]*users.User{}, nil)

	resp, err := service.GetFeed(ctx, &postsapi.FeedRequest{
		Type: postsapi.FeedType_FEED_TYPE_UNSPECIFIED,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Posts, 1)
}

func TestService_GetFeed_Subscriptions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mockPosts, mockUsers, mockSubscriptions := initService(ctrl)

	ctx := contextWithUserID("user-123")
	userID := "user-123"

	postsList := []posts.Post{
		{ID: "post-1", Text: "Post 1", UserID: "user-1", CreatedAt: time.Now()},
		{ID: "post-3", Text: "Post 3", UserID: "user-123", CreatedAt: time.Now()},
	}
	mockPosts.EXPECT().
		GetByUserIDs(ctx, []string{"user-1", "user-123"}).
		Return(postsList, nil).
		Times(1)

	mockUsers.EXPECT().
		GetByIDs(ctx, []string{"user-1", "user-123"}).
		Return(map[string]*users.User{}, nil).
		Times(1)

	mockSubscriptions.EXPECT().
		GetFollowing(ctx, userID).
		Return([]string{"user-1"}, nil).
		Times(1)

	resp, err := service.GetFeed(ctx, &postsapi.FeedRequest{
		Type: postsapi.FeedType_FEED_TYPE_SUBSCRIPTIONS,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Posts, 2)
}

func TestService_GetFeed_Subscriptions_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, _, _, _ := initService(ctrl)

	_, err := service.GetFeed(context.Background(), &postsapi.FeedRequest{
		Type: postsapi.FeedType_FEED_TYPE_SUBSCRIPTIONS,
	})
	api.RequireError(t, err, "auth_required")
}
