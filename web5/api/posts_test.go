package api

import (
	"github.com/materkov/meme9/web5/pkg/testutils"
	"github.com/materkov/meme9/web5/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandlePostsId(t *testing.T) {
	testutils.SetupRedis(t)

	require.NoError(t, store.NodeSave(9417, store.Post{
		ID:     9417,
		Date:   1666973391,
		Text:   "test post",
		UserID: 816,
	}))
	require.NoError(t, store.NodeSave(9418, store.Post{
		IsDeleted: true,
		Text:      "test post",
		UserID:    816,
	}))

	t.Run("normal post", func(t *testing.T) {
		results := handlePostsId(testutils.PrepareContext(), 15, "/posts/9417")
		post := results[0].(Post)

		require.Equal(t, "9417", post.ID)
		require.Equal(t, "2022-10-28T16:09:51Z", post.Date)
		require.Equal(t, "test post", post.Text)
		require.Equal(t, "816", post.UserID)
		require.False(t, post.IsDeleted)
		require.False(t, post.CanDelete)
	})

	t.Run("viewer is author", func(t *testing.T) {
		results := handlePostsId(testutils.PrepareContext(), 816, "/posts/9417")
		post := results[0].(Post)
		require.True(t, post.CanDelete)
	})

	t.Run("deleted post", func(t *testing.T) {
		results := handlePostsId(testutils.PrepareContext(), 816, "/posts/9418")
		post := results[0].(Post)
		require.True(t, post.IsDeleted)
		require.Equal(t, "9418", post.ID)

		require.False(t, post.CanDelete)
		require.Empty(t, post.Text)
	})

}
