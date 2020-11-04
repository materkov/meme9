package handlers

import (
	"testing"

	"github.com/materkov/meme9/api/api"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
	"github.com/stretchr/testify/require"
)

func TestUserPage_Handle(t *testing.T) {
	testStore, closer := InitTestStore(t)
	defer closer()

	require.NoError(t, testStore.AddUser(&store.User{
		ID:   44,
		Name: "test user",
	}))

	handler := UserPage{Store: testStore}

	t.Run("normal", func(t *testing.T) {
		resp, err := handler.Handle(&api.Viewer{User: &store.User{ID: 44}}, &login.UserPageRequest{UserId: "44"})
		require.NoError(t, err)
		require.Equal(t, resp.Id, "44")
		require.Equal(t, resp.LastPostId, "2")
		require.Equal(t, resp.LastPostUrl, "/posts/2")
		require.Equal(t, resp.Name, "test user")
		require.Equal(t, resp.HeaderRenderer.CurrentUserId, "44")
	})

	t.Run("no auth", func(t *testing.T) {
		resp, err := handler.Handle(&api.Viewer{}, &login.UserPageRequest{UserId: "44"})
		require.NoError(t, err)
		require.Equal(t, resp.Id, "44")
		require.Equal(t, resp.HeaderRenderer.CurrentUserId, "")
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := handler.Handle(&api.Viewer{}, &login.UserPageRequest{UserId: "44111"})
		require.Equal(t, err.(*api.Error).Code, "USER_NOT_FOUND")
	})
}
