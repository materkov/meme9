package posts

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
		collection := client.Database("meme9_test").Collection("posts")
		err = collection.Drop(ctx)
		require.NoError(t, err)
		_ = client.Disconnect(ctx)
	}

	return adapter, cleanup
}

func TestAdapter_Add(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	post := Post{
		Text:      "Test post",
		UserID:    "user-123",
		CreatedAt: time.Now(),
	}

	result, err := adapter.Add(ctx, post)
	require.NoError(t, err)
	require.NotEmpty(t, result.ID)
	require.Equal(t, post.Text, result.Text)
	require.Equal(t, post.UserID, result.UserID)
}

func TestAdapter_GetByID_Success(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	post := Post{
		Text:      "Test post",
		UserID:    "user-123",
		CreatedAt: time.Now(),
	}

	created, err := adapter.Add(ctx, post)
	require.NoError(t, err)

	retrieved, err := adapter.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, retrieved.ID)
	require.Equal(t, post.Text, retrieved.Text)
	require.Equal(t, post.UserID, retrieved.UserID)
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

func TestAdapter_GetByUserIDs(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	userID1 := "user-123"
	userID2 := "user-456"

	// Create posts for user1
	post1 := Post{Text: "Post 1", UserID: userID1, CreatedAt: time.Now()}
	post2 := Post{Text: "Post 2", UserID: userID1, CreatedAt: time.Now()}
	_, err := adapter.Add(ctx, post1)
	require.NoError(t, err)
	_, err = adapter.Add(ctx, post2)
	require.NoError(t, err)

	// Create posts for user2
	post3 := Post{Text: "Post 3", UserID: userID2, CreatedAt: time.Now()}
	_, err = adapter.Add(ctx, post3)
	require.NoError(t, err)

	// Create post for user3 (not in query)
	post4 := Post{Text: "Post 4", UserID: "user-789", CreatedAt: time.Now()}
	_, err = adapter.Add(ctx, post4)
	require.NoError(t, err)

	// Get posts for user1 and user2
	posts, err := adapter.GetByUserIDs(ctx, []string{userID1, userID2})
	require.NoError(t, err)
	require.Len(t, posts, 3)

	// Verify all posts belong to the requested users
	userIDs := make(map[string]bool)
	for _, post := range posts {
		userIDs[post.UserID] = true
		require.Contains(t, []string{userID1, userID2}, post.UserID)
	}
	require.True(t, userIDs[userID1])
	require.True(t, userIDs[userID2])
}

func TestAdapter_GetByUserIDs_Empty(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	posts, err := adapter.GetByUserIDs(ctx, []string{})
	require.NoError(t, err)
	require.Empty(t, posts)
}

func TestAdapter_GetByUserIDs_ExcludesDeleted(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()
	userID1 := "user-123"
	userID2 := "user-456"

	// Create posts for user1
	post1 := Post{Text: "Post 1", UserID: userID1, CreatedAt: time.Now()}
	post2 := Post{Text: "Post 2", UserID: userID1, CreatedAt: time.Now()}
	created1, err := adapter.Add(ctx, post1)
	require.NoError(t, err)
	created2, err := adapter.Add(ctx, post2)
	require.NoError(t, err)

	// Create posts for user2
	post3 := Post{Text: "Post 3", UserID: userID2, CreatedAt: time.Now()}
	created3, err := adapter.Add(ctx, post3)
	require.NoError(t, err)

	// Mark post2 as deleted
	err = adapter.MarkAsDeleted(ctx, created2.ID)
	require.NoError(t, err)

	// Get posts for user1 and user2 - should exclude deleted post
	posts, err := adapter.GetByUserIDs(ctx, []string{userID1, userID2})
	require.NoError(t, err)
	require.Len(t, posts, 2)

	// Verify deleted post is not included
	postIDs := make(map[string]bool)
	for _, post := range posts {
		postIDs[post.ID] = true
		require.False(t, post.Deleted, "All returned posts should not be deleted")
		require.Contains(t, []string{userID1, userID2}, post.UserID)
	}
	require.True(t, postIDs[created1.ID])
	require.True(t, postIDs[created3.ID])
	require.False(t, postIDs[created2.ID], "Deleted post should not be in results")
}

func TestAdapter_GetAll(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple posts
	post1 := Post{Text: "Post 1", UserID: "user-1", CreatedAt: time.Now()}
	post2 := Post{Text: "Post 2", UserID: "user-2", CreatedAt: time.Now()}
	post3 := Post{Text: "Post 3", UserID: "user-3", CreatedAt: time.Now()}

	_, err := adapter.Add(ctx, post1)
	require.NoError(t, err)
	_, err = adapter.Add(ctx, post2)
	require.NoError(t, err)
	_, err = adapter.Add(ctx, post3)
	require.NoError(t, err)

	// Get all posts
	posts, err := adapter.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, posts, 3)

	// Verify posts are sorted by _id descending (newest first)
	require.Equal(t, "Post 3", posts[0].Text)
	require.Equal(t, "Post 2", posts[1].Text)
	require.Equal(t, "Post 1", posts[2].Text)
}

func TestAdapter_GetAll_Empty(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	posts, err := adapter.GetAll(ctx)
	require.NoError(t, err)
	require.Empty(t, posts)
}

func TestAdapter_GetAll_ExcludesDeleted(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple posts
	post1 := Post{Text: "Post 1", UserID: "user-1", CreatedAt: time.Now()}
	post2 := Post{Text: "Post 2", UserID: "user-2", CreatedAt: time.Now()}
	post3 := Post{Text: "Post 3", UserID: "user-3", CreatedAt: time.Now()}

	created1, err := adapter.Add(ctx, post1)
	require.NoError(t, err)
	created2, err := adapter.Add(ctx, post2)
	require.NoError(t, err)
	created3, err := adapter.Add(ctx, post3)
	require.NoError(t, err)

	// Mark post2 as deleted
	err = adapter.MarkAsDeleted(ctx, created2.ID)
	require.NoError(t, err)

	// Get all posts - should exclude deleted post
	posts, err := adapter.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, posts, 2)

	// Verify deleted post is not included
	postIDs := make(map[string]bool)
	for _, post := range posts {
		postIDs[post.ID] = true
		require.False(t, post.Deleted, "All returned posts should not be deleted")
	}
	require.True(t, postIDs[created1.ID])
	require.True(t, postIDs[created3.ID])
	require.False(t, postIDs[created2.ID], "Deleted post should not be in results")
}

func TestAdapter_MarkAsDeleted(t *testing.T) {
	adapter, cleanup := setupTestAdapter(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		post := Post{
			Text:      "Test post",
			UserID:    "user-123",
			CreatedAt: time.Now(),
		}
		created, err := adapter.Add(ctx, post)
		require.NoError(t, err)

		// Mark as deleted
		err = adapter.MarkAsDeleted(ctx, created.ID)
		require.NoError(t, err)

		// Verify the post is marked as deleted
		retrieved, err := adapter.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.True(t, retrieved.Deleted)
	})

	t.Run("invalid post id", func(t *testing.T) {
		invalidID := "invalid-id"
		err := adapter.MarkAsDeleted(ctx, invalidID)
		require.NoError(t, err)
	})

	t.Run("non-existent post", func(t *testing.T) {
		nonExistentID := "507f1f77bcf86cd799439011"
		err := adapter.MarkAsDeleted(ctx, nonExistentID)
		require.NoError(t, err)
	})
}
