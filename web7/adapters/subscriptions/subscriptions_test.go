package subscriptions

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

	adapter := New(client)

	// Ensure indexes
	err = adapter.EnsureIndexes(ctx)
	require.NoError(t, err)

	// Cleanup function
	cleanup := func() {
		collection := client.Database("meme9").Collection("subscriptions")
		err = collection.Drop(ctx)
		require.NoError(t, err)
		client.Disconnect(ctx)
	}

	return adapter, cleanup
}

func TestAdapter_Subscribe_Success(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID := "user2"

	err := adapter.Subscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Verify subscription exists
	isSubscribed, err := adapter.IsSubscribed(ctx, followerID, followingID)
	require.NoError(t, err)
	require.True(t, isSubscribed)
}

func TestAdapter_Subscribe_SelfSubscription(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	userID := "user1"

	err := adapter.Subscribe(ctx, userID, userID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot subscribe to yourself")
}

func TestAdapter_Subscribe_Duplicate(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID := "user2"

	// First subscription
	err := adapter.Subscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Duplicate subscription (should be idempotent)
	err = adapter.Subscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Verify still subscribed
	isSubscribed, err := adapter.IsSubscribed(ctx, followerID, followingID)
	require.NoError(t, err)
	require.True(t, isSubscribed)
}

func TestAdapter_Unsubscribe_Success(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID := "user2"

	// Subscribe first
	err := adapter.Subscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Unsubscribe
	err = adapter.Unsubscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Verify subscription removed
	isSubscribed, err := adapter.IsSubscribed(ctx, followerID, followingID)
	require.NoError(t, err)
	require.False(t, isSubscribed)
}

func TestAdapter_Unsubscribe_NotSubscribed(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID := "user2"

	// Unsubscribe without subscribing first (should not error)
	err := adapter.Unsubscribe(ctx, followerID, followingID)
	require.NoError(t, err)
}

func TestAdapter_GetFollowing_Empty(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"

	followingIDs, err := adapter.GetFollowing(ctx, followerID)
	require.NoError(t, err)
	require.Empty(t, followingIDs)
}

func TestAdapter_GetFollowing_Multiple(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID1 := "user2"
	followingID2 := "user3"
	followingID3 := "user4"

	// Subscribe to multiple users
	err := adapter.Subscribe(ctx, followerID, followingID1)
	require.NoError(t, err)
	err = adapter.Subscribe(ctx, followerID, followingID2)
	require.NoError(t, err)
	err = adapter.Subscribe(ctx, followerID, followingID3)
	require.NoError(t, err)

	// Get following list
	followingIDs, err := adapter.GetFollowing(ctx, followerID)
	require.NoError(t, err)
	require.Len(t, followingIDs, 3)
	require.Contains(t, followingIDs, followingID1)
	require.Contains(t, followingIDs, followingID2)
	require.Contains(t, followingIDs, followingID3)
}

func TestAdapter_GetFollowing_OnlyOwnSubscriptions(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID1 := "user1"
	followerID2 := "user2"
	followingID := "user3"

	// User1 subscribes to user3
	err := adapter.Subscribe(ctx, followerID1, followingID)
	require.NoError(t, err)

	// User2 subscribes to user3
	err = adapter.Subscribe(ctx, followerID2, followingID)
	require.NoError(t, err)

	// User1 should only see their own subscription
	followingIDs, err := adapter.GetFollowing(ctx, followerID1)
	require.NoError(t, err)
	require.Len(t, followingIDs, 1)
	require.Contains(t, followingIDs, followingID)
	require.NotContains(t, followingIDs, followerID2)
}

func TestAdapter_IsSubscribed_True(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID := "user2"

	// Subscribe
	err := adapter.Subscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Check subscription
	isSubscribed, err := adapter.IsSubscribed(ctx, followerID, followingID)
	require.NoError(t, err)
	require.True(t, isSubscribed)
}

func TestAdapter_IsSubscribed_False(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID := "user2"

	// Check subscription without subscribing
	isSubscribed, err := adapter.IsSubscribed(ctx, followerID, followingID)
	require.NoError(t, err)
	require.False(t, isSubscribed)
}

func TestAdapter_IsSubscribed_AfterUnsubscribe(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID := "user2"

	// Subscribe
	err := adapter.Subscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Verify subscribed
	isSubscribed, err := adapter.IsSubscribed(ctx, followerID, followingID)
	require.NoError(t, err)
	require.True(t, isSubscribed)

	// Unsubscribe
	err = adapter.Unsubscribe(ctx, followerID, followingID)
	require.NoError(t, err)

	// Verify not subscribed
	isSubscribed, err = adapter.IsSubscribed(ctx, followerID, followingID)
	require.NoError(t, err)
	require.False(t, isSubscribed)
}

func TestAdapter_EnsureIndexes(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017/meme9?authSource=admin"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	adapter := New(client)

	// First call should succeed
	err = adapter.EnsureIndexes(ctx)
	require.NoError(t, err)

	// Second call should also succeed (idempotent)
	err = adapter.EnsureIndexes(ctx)
	require.NoError(t, err)
}

func TestAdapter_Subscribe_MultipleFollowers(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID1 := "user1"
	followerID2 := "user2"
	followerID3 := "user3"
	followingID := "user4"

	// Multiple users subscribe to the same user
	err := adapter.Subscribe(ctx, followerID1, followingID)
	require.NoError(t, err)
	err = adapter.Subscribe(ctx, followerID2, followingID)
	require.NoError(t, err)
	err = adapter.Subscribe(ctx, followerID3, followingID)
	require.NoError(t, err)

	// Verify all subscriptions exist
	isSubscribed1, err := adapter.IsSubscribed(ctx, followerID1, followingID)
	require.NoError(t, err)
	require.True(t, isSubscribed1)

	isSubscribed2, err := adapter.IsSubscribed(ctx, followerID2, followingID)
	require.NoError(t, err)
	require.True(t, isSubscribed2)

	isSubscribed3, err := adapter.IsSubscribed(ctx, followerID3, followingID)
	require.NoError(t, err)
	require.True(t, isSubscribed3)
}

func TestAdapter_GetFollowing_Order(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	followerID := "user1"
	followingID1 := "user2"
	followingID2 := "user3"
	followingID3 := "user4"

	// Subscribe in specific order
	err := adapter.Subscribe(ctx, followerID, followingID1)
	require.NoError(t, err)
	err = adapter.Subscribe(ctx, followerID, followingID3)
	require.NoError(t, err)
	err = adapter.Subscribe(ctx, followerID, followingID2)
	require.NoError(t, err)

	// Get following list
	followingIDs, err := adapter.GetFollowing(ctx, followerID)
	require.NoError(t, err)
	require.Len(t, followingIDs, 3)
	require.Contains(t, followingIDs, followingID1)
	require.Contains(t, followingIDs, followingID2)
	require.Contains(t, followingIDs, followingID3)
}
