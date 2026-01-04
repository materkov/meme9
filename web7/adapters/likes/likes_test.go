package likes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestAdapter(t *testing.T) (*Adapter, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017/meme9?authSource=admin"))
	require.NoError(t, err)

	adapter := New(client, "meme9_test")

	// Ensure indexes
	err = adapter.EnsureIndexes(ctx)
	require.NoError(t, err)

	// Cleanup function
	cleanup := func() {
		collection := client.Database("meme9_test").Collection("likes")
		err = collection.Drop(ctx)
		require.NoError(t, err)
		_ = client.Disconnect(ctx)
	}

	return adapter, cleanup
}

func TestAdapter_EnsureIndexes(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017/meme9?authSource=admin"))
	require.NoError(t, err)
	defer func() { _ = client.Disconnect(ctx) }()

	adapter := New(client, "meme9_test")

	// First call should succeed
	err = adapter.EnsureIndexes(ctx)
	require.NoError(t, err)

	// Second call should also succeed (idempotent)
	err = adapter.EnsureIndexes(ctx)
	require.NoError(t, err)
}

func TestAdapter_Like(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-123"
	userID := "user-456"

	// First like
	err := adapter.Like(ctx, postID, userID)
	require.NoError(t, err)

	// Duplicate like (should be idempotent)
	err = adapter.Like(ctx, postID, userID)
	require.NoError(t, err)

	// Verify still liked
	isLiked, err := adapter.IsLiked(ctx, postID, userID)
	require.NoError(t, err)
	require.True(t, isLiked)

	// Unlike
	err = adapter.Unlike(ctx, postID, userID)
	require.NoError(t, err)

	// Verify like removed
	isLiked, err = adapter.IsLiked(ctx, postID, userID)
	require.NoError(t, err)
	require.False(t, isLiked)
}

func TestAdapter_Unlike_NotLiked(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-123"
	userID := "user-456"

	// Unlike without liking first (should not error)
	err := adapter.Unlike(ctx, postID, userID)
	require.NoError(t, err)
}

func TestAdapter_GetLikesCounts_SinglePost(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-1"
	userID1 := "user-1"
	userID2 := "user-2"

	// Like post
	err := adapter.Like(ctx, postID, userID1)
	require.NoError(t, err)
	err = adapter.Like(ctx, postID, userID2)
	require.NoError(t, err)

	// Get counts
	counts, err := adapter.GetLikesCounts(ctx, []string{postID, "post-2"})
	require.NoError(t, err)
	require.Equal(t, 2, counts[postID])
	require.Equal(t, 0, counts["post-2"])
}

func TestAdapter_GetLikedByUser_SinglePost_Liked(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-1"
	userID := "user-456"

	// Like
	err := adapter.Like(ctx, postID, userID)
	require.NoError(t, err)

	// Get liked by user
	liked, err := adapter.GetLikedByUser(ctx, userID, []string{postID, "post-2"})
	require.NoError(t, err)
	require.True(t, liked["post-1"])
	require.False(t, liked["post-2"])
}
