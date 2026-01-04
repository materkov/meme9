package likes

import (
	"context"
	"fmt"
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
	require.Equal(t, int32(2), counts[postID])
	require.Equal(t, int32(0), counts["post-2"])
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

func TestAdapter_GetLikers_Empty(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-1"

	// Get likers for post with no likes
	userIDs, pageToken, err := adapter.GetLikers(ctx, postID, "", 10)
	require.NoError(t, err)
	require.Empty(t, userIDs)
	require.Empty(t, pageToken)
}

func TestAdapter_GetLikers_SingleLiker(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-1"
	userID := "user-1"

	// Like post
	err := adapter.Like(ctx, postID, userID)
	require.NoError(t, err)

	// Get likers
	userIDs, pageToken, err := adapter.GetLikers(ctx, postID, "", 10)
	require.NoError(t, err)
	require.Equal(t, []string{userID}, userIDs)
	require.Empty(t, pageToken) // No more results
}

func TestAdapter_GetLikers_MultipleLikers(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-1"
	userID1 := "user-1"
	userID2 := "user-2"
	userID3 := "user-3"

	// Like post with multiple users
	err := adapter.Like(ctx, postID, userID1)
	require.NoError(t, err)
	err = adapter.Like(ctx, postID, userID2)
	require.NoError(t, err)
	err = adapter.Like(ctx, postID, userID3)
	require.NoError(t, err)

	// Get all likers
	userIDs, pageToken, err := adapter.GetLikers(ctx, postID, "", 10)
	require.NoError(t, err)
	require.Len(t, userIDs, 3)
	require.Contains(t, userIDs, userID1)
	require.Contains(t, userIDs, userID2)
	require.Contains(t, userIDs, userID3)
	require.Empty(t, pageToken) // No more results
}

func TestAdapter_GetLikers_WithPagination_MoreResults(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-1"

	// Create 5 likes
	for i := 1; i <= 5; i++ {
		err := adapter.Like(ctx, postID, fmt.Sprintf("user-%d", i))
		require.NoError(t, err)
	}

	// Get first page with count 2
	userIDs, pageToken, err := adapter.GetLikers(ctx, postID, "", 2)
	require.NoError(t, err)
	require.Len(t, userIDs, 2)
	require.NotEmpty(t, pageToken) // More results available

	// Get second page using pageToken
	userIDs2, pageToken2, err := adapter.GetLikers(ctx, postID, pageToken, 2)
	require.NoError(t, err)
	require.Len(t, userIDs2, 2)
	require.NotEmpty(t, pageToken2) // More results available

	// Get third page
	userIDs3, pageToken3, err := adapter.GetLikers(ctx, postID, pageToken2, 2)
	require.NoError(t, err)
	require.Len(t, userIDs3, 1)
	require.Empty(t, pageToken3) // No more results

	// Verify all user IDs are unique across pages
	allUserIDs := append(append(userIDs, userIDs2...), userIDs3...)
	require.Len(t, allUserIDs, 5)
}

func TestAdapter_GetLikers_WithPagination_ExactCount(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID := "post-1"

	// Create exactly 3 likes
	for i := 1; i <= 3; i++ {
		err := adapter.Like(ctx, postID, fmt.Sprintf("user-%d", i))
		require.NoError(t, err)
	}

	// Get with count 3 (exact match)
	userIDs, pageToken, err := adapter.GetLikers(ctx, postID, "", 3)
	require.NoError(t, err)
	require.Len(t, userIDs, 3)
	require.Empty(t, pageToken) // No more results
}

func TestAdapter_GetLikers_IsolatedByPost(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	postID1 := "post-1"
	postID2 := "post-2"
	userID1 := "user-1"
	userID2 := "user-2"

	// Like different posts
	err := adapter.Like(ctx, postID1, userID1)
	require.NoError(t, err)
	err = adapter.Like(ctx, postID2, userID2)
	require.NoError(t, err)

	// Get likers for post-1 should only return user-1
	userIDs, pageToken, err := adapter.GetLikers(ctx, postID1, "", 10)
	require.NoError(t, err)
	require.Equal(t, []string{userID1}, userIDs)
	require.Empty(t, pageToken)

	// Get likers for post-2 should only return user-2
	userIDs2, pageToken2, err := adapter.GetLikers(ctx, postID2, "", 10)
	require.NoError(t, err)
	require.Equal(t, []string{userID2}, userIDs2)
	require.Empty(t, pageToken2)
}
