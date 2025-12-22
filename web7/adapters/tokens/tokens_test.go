package tokens

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

	// Cleanup function
	cleanup := func() {
		collection := client.Database("meme9_test").Collection("tokens")
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
	token := Token{
		Token:     "test-token-123",
		UserID:    "user-123",
		CreatedAt: time.Now(),
	}

	tokenID, err := adapter.Create(ctx, token)
	require.NoError(t, err)
	require.NotEmpty(t, tokenID)
}

func TestAdapter_GetByValue_Success(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	token := Token{
		Token:     "test-token-123",
		UserID:    "user-123",
		CreatedAt: time.Now(),
	}

	_, err := adapter.Create(ctx, token)
	require.NoError(t, err)

	retrieved, err := adapter.GetByValue(ctx, token.Token)
	require.NoError(t, err)
	require.Equal(t, token.Token, retrieved.Token)
	require.Equal(t, token.UserID, retrieved.UserID)
	require.NotEmpty(t, retrieved.ID)
}

func TestAdapter_GetByValue_NotFound(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	_, err := adapter.GetByValue(ctx, "nonexistent-token")
	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)
}

func TestAdapter_GetByValue_MultipleTokens(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple tokens for different users
	token1 := Token{Token: "token-1", UserID: "user-1", CreatedAt: time.Now()}
	token2 := Token{Token: "token-2", UserID: "user-2", CreatedAt: time.Now()}
	token3 := Token{Token: "token-3", UserID: "user-3", CreatedAt: time.Now()}

	_, err := adapter.Create(ctx, token1)
	require.NoError(t, err)
	_, err = adapter.Create(ctx, token2)
	require.NoError(t, err)
	_, err = adapter.Create(ctx, token3)
	require.NoError(t, err)

	// Retrieve each token
	retrieved1, err := adapter.GetByValue(ctx, "token-1")
	require.NoError(t, err)
	require.Equal(t, "user-1", retrieved1.UserID)

	retrieved2, err := adapter.GetByValue(ctx, "token-2")
	require.NoError(t, err)
	require.Equal(t, "user-2", retrieved2.UserID)

	retrieved3, err := adapter.GetByValue(ctx, "token-3")
	require.NoError(t, err)
	require.Equal(t, "user-3", retrieved3.UserID)
}

func TestAdapter_GetByValue_SameUserDifferentTokens(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	userID := "user-123"

	// Create multiple tokens for the same user
	token1 := Token{Token: "token-1", UserID: userID, CreatedAt: time.Now()}
	token2 := Token{Token: "token-2", UserID: userID, CreatedAt: time.Now()}

	_, err := adapter.Create(ctx, token1)
	require.NoError(t, err)
	_, err = adapter.Create(ctx, token2)
	require.NoError(t, err)

	// Retrieve each token
	retrieved1, err := adapter.GetByValue(ctx, "token-1")
	require.NoError(t, err)
	require.Equal(t, userID, retrieved1.UserID)

	retrieved2, err := adapter.GetByValue(ctx, "token-2")
	require.NoError(t, err)
	require.Equal(t, userID, retrieved2.UserID)

	// Verify they are different tokens
	require.NotEqual(t, retrieved1.ID, retrieved2.ID)
}
