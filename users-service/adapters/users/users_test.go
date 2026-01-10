//go:build integration

package users

import (
	"context"
	"testing"
	"time"

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
		err := client.Database("meme9_test").Collection("users").Drop(ctx)
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

func TestAdapter_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("get existing user by ID", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		user := User{
			Username:     "testuser",
			PasswordHash: "hash",
			CreatedAt:    time.Now(),
		}

		userID, err := adapter.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := adapter.GetByID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, userID, retrieved.ID)
		require.Equal(t, "testuser", retrieved.Username)
	})

	t.Run("get non-existent user by ID", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		nonExistentID := "507f1f77bcf86cd799439011"
		_, err := adapter.GetByID(ctx, nonExistentID)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("get user with invalid ID format", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := adapter.GetByID(ctx, "invalid-id")
		require.Equal(t, ErrNotFound, err)
	})
}

func TestAdapter_GetByIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("get multiple users", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		user1 := User{
			Username:     "user1",
			PasswordHash: "hash1",
			CreatedAt:    time.Now(),
		}
		user2 := User{
			Username:     "user2",
			PasswordHash: "hash2",
			CreatedAt:    time.Now(),
		}

		userID1, err := adapter.Create(ctx, user1)
		require.NoError(t, err)
		userID2, err := adapter.Create(ctx, user2)
		require.NoError(t, err)

		users, err := adapter.GetByIDs(ctx, []string{userID1, userID2})
		require.NoError(t, err)
		require.Len(t, users, 2)
		require.Equal(t, "user1", users[userID1].Username)
		require.Equal(t, "user2", users[userID2].Username)
	})

	t.Run("get empty list", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		users, err := adapter.GetByIDs(ctx, []string{})
		require.NoError(t, err)
		require.Empty(t, users)
	})

	t.Run("get with non-existent IDs", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		nonExistentID := "507f1f77bcf86cd799439011"
		users, err := adapter.GetByIDs(ctx, []string{nonExistentID})
		require.NoError(t, err)
		require.Empty(t, users)
	})

	t.Run("get with invalid ID format", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		users, err := adapter.GetByIDs(ctx, []string{"invalid-id"})
		require.NoError(t, err)
		require.Empty(t, users)
	})
}

func TestAdapter_UpdateAvatar(t *testing.T) {
	ctx := context.Background()

	t.Run("update avatar for existing user", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		user := User{
			Username:     "avataruser",
			PasswordHash: "hash",
			CreatedAt:    time.Now(),
		}

		userID, err := adapter.Create(ctx, user)
		require.NoError(t, err)

		avatarURL := "https://example.com/new-avatar.jpg"
		err = adapter.UpdateAvatar(ctx, userID, avatarURL)
		require.NoError(t, err)

		retrieved, err := adapter.GetByID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, avatarURL, retrieved.AvatarURL)
	})

	t.Run("update avatar for non-existent user", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		nonExistentID := "507f1f77bcf86cd799439011"
		err := adapter.UpdateAvatar(ctx, nonExistentID, "https://example.com/avatar.jpg")
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("update avatar with invalid ID format", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()

		err := adapter.UpdateAvatar(ctx, "invalid-id", "https://example.com/avatar.jpg")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid user ID")
	})
}
