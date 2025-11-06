package posts

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/services/posts/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_CreatePost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedPost := &posts.Post{
		ID:        "test-id-123",
		Text:      "Test post text",
		UserID:    "user-123",
		CreatedAt: time.Now(),
	}

	mockAdapter := mocks.NewMockPostsAdapter(ctrl)
	mockAdapter.EXPECT().
		Add(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, post posts.Post) (*posts.Post, error) {
			// Verify the post data passed to adapter
			require.Equal(t, "Test post text", post.Text)
			require.Equal(t, "user-123", post.UserID)
			require.False(t, post.CreatedAt.IsZero(), "CreatedAt should be set")
			require.WithinDuration(t, time.Now(), post.CreatedAt, time.Second, "CreatedAt should be approximately now")

			return expectedPost, nil
		})

	service := New(mockAdapter)
	result, err := service.CreatePost(ctx, "Test post text", "user-123")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedPost.ID, result.ID)
	require.Equal(t, expectedPost.Text, result.Text)
	require.Equal(t, expectedPost.UserID, result.UserID)
}

func TestService_CreatePost_AdapterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedError := errors.New("database connection failed")

	mockAdapter := mocks.NewMockPostsAdapter(ctrl)
	mockAdapter.EXPECT().
		Add(ctx, gomock.Any()).
		Return(nil, expectedError)

	service := New(mockAdapter)
	result, err := service.CreatePost(ctx, "Test post text", "user-123")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "failed to create post")
	require.ErrorIs(t, err, expectedError)
}
