package handlers

import (
	"testing"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
	"github.com/stretchr/testify/require"
)

func TestPostPage_Handle(t *testing.T) {
	testStore, closer := InitTestStore(t)
	defer closer()

	require.NoError(t, testStore.AddPost(&store.Post{
		ID:        581,
		Text:      "test text",
		UserID:    45,
		Date:      1604488449,
		UserAgent: "test user agent",
	}))

	handler := PostPage{Store: testStore}

	t.Run("Normal", func(t *testing.T) {
		viewer := &api.Viewer{User: &store.User{
			ID: 56, Name: "test user",
		}}

		r, err := handler.Handle(viewer, &pb.PostPageRequest{PostId: "581"})
		require.NoError(t, err)
		require.NotNil(t, r)
		require.Equal(t, r.Id, "581")
		require.Equal(t, r.Text, "test text")
		require.Equal(t, r.UserId, "45")
		require.Equal(t, r.UserUrl, "/users/45")
		require.Equal(t, r.PostUrl, "/posts/581")
		require.Equal(t, r.HeaderRenderer.CurrentUserId, "56")
	})

	t.Run("no auth", func(t *testing.T) {
		r, err := handler.Handle(&api.Viewer{}, &pb.PostPageRequest{PostId: "581"})
		require.NoError(t, err)
		require.Equal(t, r.Id, "581")
		require.Equal(t, "", r.HeaderRenderer.CurrentUserId)
	})

	t.Run("post not found", func(t *testing.T) {
		_, err := handler.Handle(&api.Viewer{}, &pb.PostPageRequest{PostId: "5"})
		require.Equal(t, err.(*api.Error).Code, "POST_NOT_FOUND")
	})
}
