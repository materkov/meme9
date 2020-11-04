package handlers

import (
	"testing"

	"github.com/materkov/meme9/api/api"
	login "github.com/materkov/meme9/api/pb"
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
		resp := handler.Handle(viewer, &login.PostPageRequest{PostId: "581"})

		r := resp.GetPostPageRenderer()
		require.NotNil(t, r)
		require.Equal(t, r.Id, "581")
		require.Equal(t, r.Text, "test text")
		require.Equal(t, r.UserId, "45")
		require.Equal(t, r.UserUrl, "/users/45")
		require.Equal(t, r.PostUrl, "/posts/581")
		require.Equal(t, r.HeaderRenderer.CurrentUserId, "56")
	})

	t.Run("no auth", func(t *testing.T) {
		resp := handler.Handle(&api.Viewer{}, &login.PostPageRequest{PostId: "581"})

		r := resp.GetPostPageRenderer()
		require.Equal(t, r.Id, "581")
		require.Nil(t, r.HeaderRenderer)
	})

	t.Run("post not found", func(t *testing.T) {
		resp := handler.Handle(&api.Viewer{}, &login.PostPageRequest{PostId: "5"})
		require.Equal(t, resp.GetErrorRenderer().ErrorCode, "POST_NOT_FOUND")
	})
}
