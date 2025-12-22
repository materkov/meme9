package users

import (
	"context"
	"testing"
	"time"

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
		collection := client.Database("meme9_test").Collection("users")
		err = collection.Drop(ctx)
		require.NoError(t, err)
		_ = client.Disconnect(ctx)
	}

	return adapter, cleanup
}

func TestAdapter_Create(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	user := User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	userID, err := adapter.Create(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, userID)
}

func TestAdapter_Create_DuplicateUsername(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	user := User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	_, err := adapter.Create(ctx, user)
	require.NoError(t, err)

	// Try to create another user with the same username
	_, err = adapter.Create(ctx, user)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate key")
}

func TestAdapter_GetByID_Success(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	user := User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	userID, err := adapter.Create(ctx, user)
	require.NoError(t, err)

	retrieved, err := adapter.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, retrieved.ID)
	require.Equal(t, user.Username, retrieved.Username)
	require.Equal(t, user.PasswordHash, retrieved.PasswordHash)
}

func TestAdapter_GetByID_NotFound(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	// Use a valid ObjectID format but non-existent ID
	nonExistentID := "507f1f77bcf86cd799439011"

	_, err := adapter.GetByID(ctx, nonExistentID)
	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)
}

func TestAdapter_GetByID_InvalidID(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	invalidID := "invalid-id"

	_, err := adapter.GetByID(ctx, invalidID)
	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)
}

func TestAdapter_GetByUsername_Success(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	user := User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}

	_, err := adapter.Create(ctx, user)
	require.NoError(t, err)

	retrieved, err := adapter.GetByUsername(ctx, user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, retrieved.Username)
	require.Equal(t, user.PasswordHash, retrieved.PasswordHash)
}

func TestAdapter_GetByUsername_NotFound(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	_, err := adapter.GetByUsername(ctx, "nonexistent")
	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)
}

func TestAdapter_GetByIDs_Success(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple users
	user1 := User{Username: "user1", PasswordHash: "hash1", CreatedAt: time.Now()}
	user2 := User{Username: "user2", PasswordHash: "hash2", CreatedAt: time.Now()}
	user3 := User{Username: "user3", PasswordHash: "hash3", CreatedAt: time.Now()}

	id1, err := adapter.Create(ctx, user1)
	require.NoError(t, err)
	id2, err := adapter.Create(ctx, user2)
	require.NoError(t, err)
	id3, err := adapter.Create(ctx, user3)
	require.NoError(t, err)

	// Get users by IDs
	users, err := adapter.GetByIDs(ctx, []string{id1, id2, id3})
	require.NoError(t, err)
	require.Len(t, users, 3)

	require.Equal(t, "user1", users[id1].Username)
	require.Equal(t, "user2", users[id2].Username)
	require.Equal(t, "user3", users[id3].Username)
}

func TestAdapter_GetByIDs_Partial(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	// Create a user
	user1 := User{Username: "user1", PasswordHash: "hash1", CreatedAt: time.Now()}
	id1, err := adapter.Create(ctx, user1)
	require.NoError(t, err)

	// Try to get with one valid ID and one invalid ID
	users, err := adapter.GetByIDs(ctx, []string{id1, "507f1f77bcf86cd799439011"})
	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, "user1", users[id1].Username)
}

func TestAdapter_GetByIDs_Empty(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	users, err := adapter.GetByIDs(ctx, []string{})
	require.NoError(t, err)
	require.Empty(t, users)
}

func TestAdapter_GetByIDs_AllInvalid(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	users, err := adapter.GetByIDs(ctx, []string{"invalid-id-1", "invalid-id-2"})
	require.NoError(t, err)
	require.Empty(t, users)
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

	// Cleanup
	collection := client.Database("meme9").Collection("users")
	_ = collection.Drop(ctx)
}
