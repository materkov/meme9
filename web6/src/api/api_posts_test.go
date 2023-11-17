package api

import (
	"errors"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func requireAPIError(t *testing.T, err error, code string) {
	t.Helper()
	var apiErr Error
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, string(apiErr), code)
}

func TestAPI_PostsCRUD(t *testing.T) {
	api := API{}
	closer := createTestDB(t)
	defer closer()

	userID, _ := store.GlobalStore.AddObject(store.ObjTypeUser, store.User{ID: 1})
	v := Viewer{UserID: 1}

	addResp, err := api.PostsAdd(&v, &PostsAddReq{Text: "test text"})
	require.NoError(t, err)
	require.NotNil(t, addResp)
	require.NotEmpty(t, addResp.ID)

	postID := addResp.ID

	t.Run("", func(t *testing.T) {
		resp, err := api.PostsList(&v, &PostsListReq{})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, resp.Items[0].ID, postID)
	})

	t.Run("", func(t *testing.T) {
		resp, err := api.PostsListByID(&v, &PostsListByIdReq{ID: postID})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, resp.ID, postID)
	})

	t.Run("", func(t *testing.T) {
		resp, err := api.PostsListByUser(&v, &PostsListByUserReq{UserID: strconv.Itoa(userID)})
		require.NoError(t, err)
		require.Len(t, resp, 1)
		require.Equal(t, resp[0].ID, postID)
	})

	t.Run("", func(t *testing.T) {
		_, err = api.PostsDelete(&v, &PostsDeleteReq{PostID: postID})
		require.NoError(t, err)
	})

	t.Run("", func(t *testing.T) {
		_, err := api.PostsListByID(&v, &PostsListByIdReq{ID: postID})
		requireAPIError(t, err, "PostNotFound")
	})
}

func TestAPI_PostsLikes(t *testing.T) {
	api := API{}

	closer := createTestDB(t)
	defer closer()

	userID, _ := store.GlobalStore.AddObject(store.ObjTypeUser, store.User{})
	v := Viewer{UserID: userID}

	addResp, _ := api.PostsAdd(&v, &PostsAddReq{Text: "test text"})

	t.Run("like post", func(t *testing.T) {
		_, err := api.PostsLike(&v, &PostsLikeReq{
			PostID: addResp.ID,
			Action: "LIKE",
		})
		require.NoError(t, err)
	})

	t.Run("check count and flag", func(t *testing.T) {
		listResp, err := api.PostsListByID(&v, &PostsListByIdReq{ID: addResp.ID})
		require.NoError(t, err)
		require.Equal(t, 1, listResp.LikesCount)
		require.True(t, listResp.IsLiked)
	})

	t.Run("dislike post", func(t *testing.T) {
		_, err := api.PostsLike(&v, &PostsLikeReq{
			PostID: addResp.ID,
			Action: Unlike,
		})
		require.NoError(t, err)
	})

	t.Run("check again", func(t *testing.T) {
		listResp, err := api.PostsListByID(&v, &PostsListByIdReq{ID: addResp.ID})
		require.NoError(t, err)
		require.Equal(t, 0, listResp.LikesCount)
		require.False(t, listResp.IsLiked)
	})
}
