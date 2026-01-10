//go:build integration

package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*MongoAdapter, func()) {
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
		err = client.Database("meme9_test").Collection("tokens").Drop(ctx)
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

func TestCreateUser(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful user creation", func(t *testing.T) {
		user := User{
			Username:     "testuser",
			PasswordHash: "hashedpassword",
			CreatedAt:    time.Now(),
		}

		userID, err := adapter.CreateUser(ctx, user)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		// Verify user was created
		retrieved, err := adapter.GetUserByUsername(ctx, "testuser")
		require.NoError(t, err)
		require.Equal(t, "testuser", retrieved.Username)
		require.Equal(t, "hashedpassword", retrieved.PasswordHash)
		require.Equal(t, userID, retrieved.ID)
	})
}

func TestGetUserByUsername(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("get existing user", func(t *testing.T) {
		user := User{
			Username:     "existinguser",
			PasswordHash: "hash",
			CreatedAt:    time.Now(),
		}

		userID, err := adapter.CreateUser(ctx, user)
		require.NoError(t, err)

		retrieved, err := adapter.GetUserByUsername(ctx, "existinguser")
		require.NoError(t, err)
		require.Equal(t, userID, retrieved.ID)
		require.Equal(t, "existinguser", retrieved.Username)
		require.Equal(t, "hash", retrieved.PasswordHash)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		_, err := adapter.GetUserByUsername(ctx, "nonexistent")
		require.Equal(t, ErrUserNotFound, err)
	})
}

func TestGetUserByID(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Only verify that works, everything else in batch method
	nonExistentID := "507f1f77bcf86cd799439011"
	_, err := adapter.GetUserByID(ctx, nonExistentID)
	require.Equal(t, ErrUserNotFound, err)
}

func TestGetUsersByIDs(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("get multiple users", func(t *testing.T) {
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

		userID1, err := adapter.CreateUser(ctx, user1)
		require.NoError(t, err)
		userID2, err := adapter.CreateUser(ctx, user2)
		require.NoError(t, err)

		users, err := adapter.GetUsersByIDs(ctx, []string{userID1, userID2})
		require.NoError(t, err)
		require.Len(t, users, 2)
		require.Equal(t, "user1", users[userID1].Username)
		require.Equal(t, "user2", users[userID2].Username)
	})

	t.Run("get empty list", func(t *testing.T) {
		users, err := adapter.GetUsersByIDs(ctx, []string{})
		require.NoError(t, err)
		require.Empty(t, users)
	})

	t.Run("get with non-existent IDs", func(t *testing.T) {
		nonExistentID := "507f1f77bcf86cd799439011"
		users, err := adapter.GetUsersByIDs(ctx, []string{nonExistentID})
		require.NoError(t, err)
		require.Empty(t, users)
	})

	t.Run("get with invalid ID format", func(t *testing.T) {
		users, err := adapter.GetUsersByIDs(ctx, []string{"invalid-id"})
		require.NoError(t, err)
		require.Empty(t, users)
	})

	t.Run("get with mix of valid and invalid IDs", func(t *testing.T) {
		user := User{
			Username:     "validuser",
			PasswordHash: "hash",
			CreatedAt:    time.Now(),
		}

		userID, err := adapter.CreateUser(ctx, user)
		require.NoError(t, err)

		users, err := adapter.GetUsersByIDs(ctx, []string{userID, "invalid-id", "507f1f77bcf86cd799439011"})
		require.NoError(t, err)
		require.Len(t, users, 1)
		require.Equal(t, "validuser", users[userID].Username)
	})
}

func TestUpdateUserAvatar(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("update avatar for existing user", func(t *testing.T) {
		user := User{
			Username:     "avataruser",
			PasswordHash: "hash",
			CreatedAt:    time.Now(),
		}

		userID, err := adapter.CreateUser(ctx, user)
		require.NoError(t, err)

		avatarURL := "https://example.com/new-avatar.jpg"
		err = adapter.UpdateUserAvatar(ctx, userID, avatarURL)
		require.NoError(t, err)

		retrieved, err := adapter.GetUserByID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, avatarURL, retrieved.AvatarURL)
	})

	t.Run("update avatar for non-existent user", func(t *testing.T) {
		nonExistentID := "507f1f77bcf86cd799439011"
		err := adapter.UpdateUserAvatar(ctx, nonExistentID, "https://example.com/avatar.jpg")
		require.Equal(t, ErrUserNotFound, err)
	})

	t.Run("update avatar with invalid ID format", func(t *testing.T) {
		err := adapter.UpdateUserAvatar(ctx, "invalid-id", "https://example.com/avatar.jpg")
		require.Equal(t, ErrUserNotFound, err)
	})
}

func TestCreateToken(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful token creation", func(t *testing.T) {
		token := Token{
			Token:     "test-token-123",
			UserID:    "user123",
			CreatedAt: time.Now(),
		}

		tokenID, err := adapter.CreateToken(ctx, token)
		require.NoError(t, err)
		require.NotEmpty(t, tokenID)

		// Verify token was created
		retrieved, err := adapter.GetTokenByValue(ctx, "test-token-123")
		require.NoError(t, err)
		require.Equal(t, "test-token-123", retrieved.Token)
		require.Equal(t, "user123", retrieved.UserID)
		require.Equal(t, tokenID, retrieved.ID)
	})

	t.Run("create multiple tokens for same user", func(t *testing.T) {
		token1 := Token{
			Token:     "token1",
			UserID:    "user456",
			CreatedAt: time.Now(),
		}
		token2 := Token{
			Token:     "token2",
			UserID:    "user456",
			CreatedAt: time.Now(),
		}

		_, err := adapter.CreateToken(ctx, token1)
		require.NoError(t, err)

		_, err = adapter.CreateToken(ctx, token2)
		require.NoError(t, err)

		// Both tokens should exist
		retrieved1, err := adapter.GetTokenByValue(ctx, "token1")
		require.NoError(t, err)
		require.Equal(t, "user456", retrieved1.UserID)

		retrieved2, err := adapter.GetTokenByValue(ctx, "token2")
		require.NoError(t, err)
		require.Equal(t, "user456", retrieved2.UserID)
	})
}

func TestGetTokenByValue(t *testing.T) {
	adapter, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("get existing token", func(t *testing.T) {
		token := Token{
			Token:     "existing-token",
			UserID:    "user789",
			CreatedAt: time.Now(),
		}

		tokenID, err := adapter.CreateToken(ctx, token)
		require.NoError(t, err)

		retrieved, err := adapter.GetTokenByValue(ctx, "existing-token")
		require.NoError(t, err)
		require.Equal(t, tokenID, retrieved.ID)
		require.Equal(t, "existing-token", retrieved.Token)
		require.Equal(t, "user789", retrieved.UserID)
	})

	t.Run("get non-existent token", func(t *testing.T) {
		_, err := adapter.GetTokenByValue(ctx, "nonexistent-token")
		require.Equal(t, ErrTokenNotFound, err)
	})
}
