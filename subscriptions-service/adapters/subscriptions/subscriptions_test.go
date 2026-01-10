//go:build integration

package subscriptions

import (
	"context"
	"testing"

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
		err := client.Database("meme9_test").Collection("subscriptions").Drop(ctx)
		require.NoError(t, err)
	}

	return adapter, cleanup
}

func TestAdapter_EnsureIndexes(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	err := adapter.EnsureIndexes(ctx)
	require.NoError(t, err)
}

func TestAdapter_Subscribe(t *testing.T) {
	ctx := context.Background()

	t.Run("successful subscribe", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		isSubscribed, err := adapter.IsSubscribed(ctx, "user1", "user2")
		require.NoError(t, err)
		require.True(t, isSubscribed)
	})

	t.Run("idempotent subscribe", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		err = adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		isSubscribed, err := adapter.IsSubscribed(ctx, "user1", "user2")
		require.NoError(t, err)
		require.True(t, isSubscribed)
	})

	t.Run("cannot subscribe to yourself", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Subscribe(ctx, "user1", "user1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot subscribe to yourself")
	})
}

func TestAdapter_Unsubscribe(t *testing.T) {
	ctx := context.Background()

	t.Run("successful unsubscribe", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		err = adapter.Unsubscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		isSubscribed, err := adapter.IsSubscribed(ctx, "user1", "user2")
		require.NoError(t, err)
		require.False(t, isSubscribed)
	})

	t.Run("unsubscribe non-existent subscription", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Unsubscribe(ctx, "user1", "user2")
		require.NoError(t, err)
	})
}

func TestAdapter_IsSubscribed(t *testing.T) {
	ctx := context.Background()

	t.Run("is subscribed", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		isSubscribed, err := adapter.IsSubscribed(ctx, "user1", "user2")
		require.NoError(t, err)
		require.True(t, isSubscribed)
	})

	t.Run("not subscribed", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		isSubscribed, err := adapter.IsSubscribed(ctx, "user1", "user2")
		require.NoError(t, err)
		require.False(t, isSubscribed)
	})
}

func TestAdapter_GetFollowing(t *testing.T) {
	ctx := context.Background()

	t.Run("get following list", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		err = adapter.Subscribe(ctx, "user1", "user3")
		require.NoError(t, err)

		following, err := adapter.GetFollowing(ctx, "user1")
		require.NoError(t, err)
		require.Len(t, following, 2)
		require.Contains(t, following, "user2")
		require.Contains(t, following, "user3")
	})

	t.Run("empty following list", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		following, err := adapter.GetFollowing(ctx, "user1")
		require.NoError(t, err)
		require.Empty(t, following)
	})

	t.Run("get following after unsubscribe", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		err = adapter.Subscribe(ctx, "user1", "user3")
		require.NoError(t, err)

		err = adapter.Unsubscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		following, err := adapter.GetFollowing(ctx, "user1")
		require.NoError(t, err)
		require.Len(t, following, 1)
		require.Contains(t, following, "user3")
		require.NotContains(t, following, "user2")
	})
}

func TestAdapter_Integration(t *testing.T) {
	ctx := context.Background()

	t.Run("full workflow", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		// Subscribe
		err := adapter.Subscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		// Check status
		isSubscribed, err := adapter.IsSubscribed(ctx, "user1", "user2")
		require.NoError(t, err)
		require.True(t, isSubscribed)

		// Get following
		following, err := adapter.GetFollowing(ctx, "user1")
		require.NoError(t, err)
		require.Contains(t, following, "user2")

		// Unsubscribe
		err = adapter.Unsubscribe(ctx, "user1", "user2")
		require.NoError(t, err)

		// Check status again
		isSubscribed, err = adapter.IsSubscribed(ctx, "user1", "user2")
		require.NoError(t, err)
		require.False(t, isSubscribed)

		// Get following again
		following, err = adapter.GetFollowing(ctx, "user1")
		require.NoError(t, err)
		require.NotContains(t, following, "user2")
	})
}
