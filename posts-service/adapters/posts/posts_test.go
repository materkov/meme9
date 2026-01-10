//go:build integration

package posts

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

	// Cleanup function
	cleanup := func() {
		err := client.Database("meme9_test").Collection("posts").Drop(ctx)
		require.NoError(t, err)
	}

	return adapter, cleanup
}

func TestAdapter_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("get all posts", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post1 := Post{
			Text:      "First post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}
		post2 := Post{
			Text:      "Second post",
			UserID:    "user2",
			CreatedAt: time.Now(),
		}

		_, err := adapter.Add(ctx, post1)
		require.NoError(t, err)

		_, err = adapter.Add(ctx, post2)
		require.NoError(t, err)

		posts, err := adapter.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, posts, 2)
		require.Equal(t, "Second post", posts[0].Text) // Newest first
		require.Equal(t, "First post", posts[1].Text)
	})

	t.Run("exclude deleted posts", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		ctx := context.Background()
		post1 := Post{
			Text:      "Active post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}
		post2 := Post{
			Text:      "Deleted post",
			UserID:    "user2",
			CreatedAt: time.Now(),
		}

		_, err := adapter.Add(ctx, post1)
		require.NoError(t, err)

		added2, err := adapter.Add(ctx, post2)
		require.NoError(t, err)

		err = adapter.MarkAsDeleted(ctx, added2.ID)
		require.NoError(t, err)

		posts, err := adapter.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, posts, 1)
		require.Equal(t, "Active post", posts[0].Text)
	})

	t.Run("empty result", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		ctx := context.Background()
		posts, err := adapter.GetAll(ctx)
		require.NoError(t, err)
		require.Empty(t, posts)
	})
}

func TestAdapter_GetByUserIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("get posts by user IDs", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post1 := Post{
			Text:      "User1 post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}
		post2 := Post{
			Text:      "User2 post",
			UserID:    "user2",
			CreatedAt: time.Now(),
		}
		post3 := Post{
			Text:      "User3 post",
			UserID:    "user3",
			CreatedAt: time.Now(),
		}

		_, err := adapter.Add(ctx, post1)
		require.NoError(t, err)

		_, err = adapter.Add(ctx, post2)
		require.NoError(t, err)

		_, err = adapter.Add(ctx, post3)
		require.NoError(t, err)

		posts, err := adapter.GetByUserIDs(ctx, []string{"user1", "user2"})
		require.NoError(t, err)
		require.Len(t, posts, 2)
		require.Equal(t, "User2 post", posts[0].Text) // Newest first
		require.Equal(t, "User1 post", posts[1].Text)
	})

	t.Run("exclude deleted posts", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		ctx := context.Background()
		post1 := Post{
			Text:      "Active post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}
		post2 := Post{
			Text:      "Deleted post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}

		_, err := adapter.Add(ctx, post1)
		require.NoError(t, err)

		added2, err := adapter.Add(ctx, post2)
		require.NoError(t, err)

		err = adapter.MarkAsDeleted(ctx, added2.ID)
		require.NoError(t, err)

		posts, err := adapter.GetByUserIDs(ctx, []string{"user1"})
		require.NoError(t, err)
		require.Len(t, posts, 1)
		require.Equal(t, "Active post", posts[0].Text)
	})

	t.Run("empty user IDs list", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		posts, err := adapter.GetByUserIDs(ctx, []string{})
		require.NoError(t, err)
		require.Empty(t, posts)
	})

	t.Run("no matching users", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post := Post{
			Text:      "Post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}

		_, err := adapter.Add(ctx, post)
		require.NoError(t, err)

		posts, err := adapter.GetByUserIDs(ctx, []string{"nonexistent"})
		require.NoError(t, err)
		require.Empty(t, posts)
	})
}

func TestAdapter_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("get existing post", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post := Post{
			Text:      "Test post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}

		added, err := adapter.Add(ctx, post)
		require.NoError(t, err)

		retrieved, err := adapter.GetByID(ctx, added.ID)
		require.NoError(t, err)
		require.Equal(t, added.ID, retrieved.ID)
		require.Equal(t, "Test post", retrieved.Text)
		require.Equal(t, "user1", retrieved.UserID)
	})

	t.Run("get non-existent post", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		nonExistentID := "507f1f77bcf86cd799439011"
		_, err := adapter.GetByID(ctx, nonExistentID)
		require.Error(t, err)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("invalid ID format", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		_, err := adapter.GetByID(ctx, "invalid-id")
		require.Error(t, err)
		require.Equal(t, ErrNotFound, err)
	})
}

func TestAdapter_Add(t *testing.T) {
	ctx := context.Background()

	t.Run("add post", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post := Post{
			Text:      "New post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}

		added, err := adapter.Add(ctx, post)
		require.NoError(t, err)
		require.NotEmpty(t, added.ID)
		require.Equal(t, "New post", added.Text)
		require.Equal(t, "user1", added.UserID)
		require.False(t, added.Deleted)

		// Verify it was saved
		retrieved, err := adapter.GetByID(ctx, added.ID)
		require.NoError(t, err)
		require.Equal(t, "New post", retrieved.Text)
	})

	t.Run("add multiple posts", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post1 := Post{
			Text:      "Post 1",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}
		post2 := Post{
			Text:      "Post 2",
			UserID:    "user2",
			CreatedAt: time.Now(),
		}

		added1, err := adapter.Add(ctx, post1)
		require.NoError(t, err)

		added2, err := adapter.Add(ctx, post2)
		require.NoError(t, err)

		require.NotEqual(t, added1.ID, added2.ID)

		posts, err := adapter.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, posts, 2)
	})
}

func TestAdapter_MarkAsDeleted(t *testing.T) {
	ctx := context.Background()

	t.Run("mark post as deleted", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post := Post{
			Text:      "Post to delete",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}

		added, err := adapter.Add(ctx, post)
		require.NoError(t, err)

		err = adapter.MarkAsDeleted(ctx, added.ID)
		require.NoError(t, err)

		// Verify it's marked as deleted
		retrieved, err := adapter.GetByID(ctx, added.ID)
		require.NoError(t, err)
		require.True(t, retrieved.Deleted)

		// Verify it's excluded from GetAll
		posts, err := adapter.GetAll(ctx)
		require.NoError(t, err)
		require.Empty(t, posts)
	})

	t.Run("mark non-existent post", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		nonExistentID := "507f1f77bcf86cd799439011"
		err := adapter.MarkAsDeleted(ctx, nonExistentID)
		require.NoError(t, err) // Returns nil for invalid IDs
	})

	t.Run("mark with invalid ID format", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		err := adapter.MarkAsDeleted(ctx, "invalid-id")
		require.NoError(t, err) // Returns nil for invalid IDs
	})

	t.Run("mark already deleted post", func(t *testing.T) {
		adapter, cleanup := setupTestDB(t)
		defer cleanup()
		post := Post{
			Text:      "Post",
			UserID:    "user1",
			CreatedAt: time.Now(),
		}

		added, err := adapter.Add(ctx, post)
		require.NoError(t, err)

		err = adapter.MarkAsDeleted(ctx, added.ID)
		require.NoError(t, err)

		// Mark again
		err = adapter.MarkAsDeleted(ctx, added.ID)
		require.NoError(t, err)

		retrieved, err := adapter.GetByID(ctx, added.ID)
		require.NoError(t, err)
		require.True(t, retrieved.Deleted)
	})
}
