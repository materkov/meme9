//go:build integration

package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*Adapter, func()) {
	t.Helper()
	mongoURI := "mongodb://admin:password@localhost:27017/meme9?authSource=admin"

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	adapter := New(client, "meme9_test")

	err = adapter.EnsureIndexes(ctx)
	require.NoError(t, err)

	// Cleanup function
	cleanup := func() {
		err := client.Database("meme9_test").Collection("likes").Drop(ctx)
		require.NoError(t, err)
	}

	return adapter, cleanup
}

func TestEnsureIndexes(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	err := adapter.EnsureIndexes(ctx)
	require.NoError(t, err)
}

func TestLike(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful like", func(t *testing.T) {
		err := adapter.Like(ctx, "post1", "user1")
		require.NoError(t, err)

		isLiked, err := adapter.IsLiked(ctx, "post1", "user1")
		require.NoError(t, err)
		assert.True(t, isLiked)
	})

	t.Run("idempotent like", func(t *testing.T) {
		// First like
		err := adapter.Like(ctx, "post2", "user1")
		require.NoError(t, err)

		// Second like (should be idempotent)
		err = adapter.Like(ctx, "post2", "user1")
		require.NoError(t, err)

		// Verify only one like exists
		counts, err := adapter.GetLikesCounts(ctx, []string{"post2"})
		require.NoError(t, err)
		assert.Equal(t, int32(1), counts["post2"])
	})

	t.Run("multiple users can like same post", func(t *testing.T) {
		err := adapter.Like(ctx, "post3", "user1")
		require.NoError(t, err)

		err = adapter.Like(ctx, "post3", "user2")
		require.NoError(t, err)

		counts, err := adapter.GetLikesCounts(ctx, []string{"post3"})
		require.NoError(t, err)
		assert.Equal(t, int32(2), counts["post3"])
	})
}

func TestUnlike(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("unlike existing like", func(t *testing.T) {
		// First create a like
		err := adapter.Like(ctx, "post1", "user1")
		require.NoError(t, err)

		// Then unlike it
		err = adapter.Unlike(ctx, "post1", "user1")
		require.NoError(t, err)

		// Verify it's no longer liked
		isLiked, err := adapter.IsLiked(ctx, "post1", "user1")
		require.NoError(t, err)
		assert.False(t, isLiked, "Post should not be liked after Unlike()")
	})

	t.Run("unlike non-existent like", func(t *testing.T) {
		// Unlike a post that was never liked (should not error)
		err := adapter.Unlike(ctx, "post2", "user1")
		require.NoError(t, err, "Unlike() on non-existent like should not fail")
	})

	t.Run("unlike one user doesn't affect others", func(t *testing.T) {
		// Create likes from multiple users
		adapter.Like(ctx, "post3", "user1")
		adapter.Like(ctx, "post3", "user2")
		adapter.Like(ctx, "post3", "user3")

		// Unlike from one user
		err := adapter.Unlike(ctx, "post3", "user2")
		require.NoError(t, err)

		// Verify other likes still exist
		counts, err := adapter.GetLikesCounts(ctx, []string{"post3"})
		require.NoError(t, err, "GetLikesCounts() failed")
		assert.Equal(t, int32(2), counts["post3"], "Expected 2 likes after unliking one")

		// Verify user1 and user3 still have likes
		isLiked1, _ := adapter.IsLiked(ctx, "post3", "user1")
		isLiked3, _ := adapter.IsLiked(ctx, "post3", "user3")
		assert.True(t, isLiked1 && isLiked3, "Other users' likes should still exist")
	})
}

func TestIsLiked(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("returns false for non-existent like", func(t *testing.T) {
		isLiked, err := adapter.IsLiked(ctx, "post1", "user1")
		require.NoError(t, err)
		assert.False(t, isLiked, "IsLiked() should return false for non-existent like")
	})

	t.Run("returns true for existing like", func(t *testing.T) {
		err := adapter.Like(ctx, "post2", "user1")
		require.NoError(t, err)

		isLiked, err := adapter.IsLiked(ctx, "post2", "user1")
		require.NoError(t, err)
		assert.True(t, isLiked, "IsLiked() should return true for existing like")
	})

	t.Run("returns false after unlike", func(t *testing.T) {
		adapter.Like(ctx, "post3", "user1")
		adapter.Unlike(ctx, "post3", "user1")

		isLiked, err := adapter.IsLiked(ctx, "post3", "user1")
		require.NoError(t, err)
		assert.False(t, isLiked, "IsLiked() should return false after Unlike()")
	})
}

func TestGetLikesCounts(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("empty postIDs returns empty map", func(t *testing.T) {
		counts, err := adapter.GetLikesCounts(ctx, []string{})
		require.NoError(t, err)
		assert.Empty(t, counts, "Expected empty map")
	})

	t.Run("returns zero for posts with no likes", func(t *testing.T) {
		counts, err := adapter.GetLikesCounts(ctx, []string{"post1", "post2"})
		require.NoError(t, err)
		assert.Equal(t, int32(0), counts["post1"], "Expected zero for post1")
		assert.Equal(t, int32(0), counts["post2"], "Expected zero for post2")
	})

	t.Run("returns correct counts for single post", func(t *testing.T) {
		// Add multiple likes
		adapter.Like(ctx, "post1", "user1")
		adapter.Like(ctx, "post1", "user2")
		adapter.Like(ctx, "post1", "user3")

		counts, err := adapter.GetLikesCounts(ctx, []string{"post1"})
		require.NoError(t, err)
		assert.Equal(t, int32(3), counts["post1"], "Expected 3 likes")
	})

	t.Run("returns correct counts for multiple posts", func(t *testing.T) {
		// Setup likes for multiple posts
		adapter.Like(ctx, "post1", "user1")
		adapter.Like(ctx, "post1", "user2")
		adapter.Like(ctx, "post2", "user1")
		adapter.Like(ctx, "post3", "user1")
		adapter.Like(ctx, "post3", "user2")
		adapter.Like(ctx, "post3", "user3")

		counts, err := adapter.GetLikesCounts(ctx, []string{"post1", "post2", "post3", "post4"})
		require.NoError(t, err)

		assert.Equal(t, int32(2), counts["post1"], "Expected 2 likes for post1")
		assert.Equal(t, int32(1), counts["post2"], "Expected 1 like for post2")
		assert.Equal(t, int32(3), counts["post3"], "Expected 3 likes for post3")
		assert.Equal(t, int32(0), counts["post4"], "Expected 0 likes for post4")
	})
}

func TestGetLikedByUser(t *testing.T) {
	ctx := context.Background()

	t.Run("empty postIDs returns empty map", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		result, err := adapter.GetLikedByUser(ctx, "user1", []string{})
		require.NoError(t, err)
		assert.Empty(t, result, "Expected empty map")
	})

	t.Run("returns false for posts not liked by user", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		result, err := adapter.GetLikedByUser(ctx, "user1", []string{"post1", "post2"})
		require.NoError(t, err)
		assert.False(t, result["post1"], "Expected post1 to be false")
		assert.False(t, result["post2"], "Expected post2 to be false")
	})

	t.Run("returns true for posts liked by user", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		adapter.Like(ctx, "post1", "user1")
		adapter.Like(ctx, "post2", "user1")
		// post3 not liked by user1

		result, err := adapter.GetLikedByUser(ctx, "user1", []string{"post1", "post2", "post3"})
		require.NoError(t, err)

		assert.True(t, result["post1"], "Expected post1 to be liked by user1")
		assert.True(t, result["post2"], "Expected post2 to be liked by user1")
		assert.False(t, result["post3"], "Expected post3 to not be liked by user1")
	})

	t.Run("only returns likes for specified user", func(t *testing.T) {
		// user1 likes post1 and post2
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		adapter.Like(ctx, "post1", "user1")
		adapter.Like(ctx, "post2", "user1")
		// user2 likes post2 and post3
		adapter.Like(ctx, "post2", "user2")
		adapter.Like(ctx, "post3", "user2")

		result, err := adapter.GetLikedByUser(ctx, "user1", []string{"post1", "post2", "post3"})
		require.NoError(t, err)

		assert.True(t, result["post1"], "Expected post1 to be liked by user1")
		assert.True(t, result["post2"], "Expected post2 to be liked by user1")
		assert.False(t, result["post3"], "Expected post3 to not be liked by user1")
	})
}

func TestGetLikers(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("returns empty for post with no likes", func(t *testing.T) {
		likers, token, err := adapter.GetLikers(ctx, "post1", "", 10)
		require.NoError(t, err)
		assert.Empty(t, likers, "Expected empty likers")
		assert.Empty(t, token, "Expected empty token")
	})

	t.Run("returns all likers when count is sufficient", func(t *testing.T) {
		// Add multiple likes
		adapter.Like(ctx, "post1", "user1")
		adapter.Like(ctx, "post1", "user2")
		adapter.Like(ctx, "post1", "user3")

		likers, token, err := adapter.GetLikers(ctx, "post1", "", 10)
		require.NoError(t, err)
		assert.Len(t, likers, 3, "Expected 3 likers")
		assert.Empty(t, token, "Expected empty token when all results fit")

		// Verify all users are present
		userMap := make(map[string]bool)
		for _, userID := range likers {
			userMap[userID] = true
		}
		assert.True(t, userMap["user1"] && userMap["user2"] && userMap["user3"], "Expected all three users in likers list")
	})

	t.Run("returns pagination token when more results available", func(t *testing.T) {
		// Add more likes than the count limit
		adapter.Like(ctx, "post2", "user1")
		adapter.Like(ctx, "post2", "user2")
		adapter.Like(ctx, "post2", "user3")
		adapter.Like(ctx, "post2", "user4")
		adapter.Like(ctx, "post2", "user5")

		// Request only 3, should get token
		likers, token, err := adapter.GetLikers(ctx, "post2", "", 3)
		require.NoError(t, err)
		assert.Len(t, likers, 3, "Expected 3 likers")
		assert.NotEmpty(t, token, "Expected pagination token when more results available")

		// Get next page
		likers2, token2, err := adapter.GetLikers(ctx, "post2", token, 3)
		require.NoError(t, err)
		assert.Len(t, likers2, 2, "Expected 2 likers on second page")
		assert.Empty(t, token2, "Expected empty token on last page")

		// Verify no duplicates
		allLikers := append(likers, likers2...)
		userMap := make(map[string]bool)
		for _, userID := range allLikers {
			assert.False(t, userMap[userID], "Duplicate user found: %s", userID)
			userMap[userID] = true
		}
	})

	t.Run("pagination with page token", func(t *testing.T) {
		adapter.Like(ctx, "post3", "user1")
		adapter.Like(ctx, "post3", "user2")
		adapter.Like(ctx, "post3", "user3")

		// Get first page
		likers1, token1, err := adapter.GetLikers(ctx, "post3", "", 2)
		require.NoError(t, err)
		assert.Len(t, likers1, 2, "Expected 2 likers on first page")
		assert.NotEmpty(t, token1, "Expected pagination token")

		// Get second page
		likers2, token2, err := adapter.GetLikers(ctx, "post3", token1, 2)
		require.NoError(t, err)
		assert.Len(t, likers2, 1, "Expected 1 liker on second page")
		assert.Empty(t, token2, "Expected empty token on last page")
	})

	t.Run("invalid page token is ignored", func(t *testing.T) {
		adapter.Like(ctx, "post4", "user1")

		// Use invalid token
		likers, token, err := adapter.GetLikers(ctx, "post4", "invalid_token", 10)
		require.NoError(t, err)
		// Should still return results (invalid token is ignored)
		assert.GreaterOrEqual(t, len(likers), 1, "Expected results even with invalid token")
		assert.Empty(t, token, "Expected empty token when all results fit")
	})
}
