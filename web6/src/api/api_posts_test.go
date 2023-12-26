package api

import (
	"context"
	"errors"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
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

	user := store.User{}
	_ = store2.GlobalStore.Users.Add(&user)
	v := Viewer{UserID: user.ID}

	addResp, err := api.PostsAdd(context.Background(), &v, &PostsAddReq{Text: "test text"})
	require.NoError(t, err)
	require.NotNil(t, addResp)
	require.NotEmpty(t, addResp.ID)

	postID := addResp.ID

	t.Run("", func(t *testing.T) {
		resp, err := api.PostsList(context.Background(), &v, &PostsListReq{})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, resp.Items[0].ID, postID)
	})

	t.Run("", func(t *testing.T) {
		resp, err := api.PostsListByID(context.Background(), &v, &PostsListByIdReq{ID: postID})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, resp.ID, postID)
	})

	t.Run("", func(t *testing.T) {
		resp, err := api.PostsList(context.Background(), &v, &PostsListReq{ByUserID: strconv.Itoa(user.ID)})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, resp.Items[0].ID, postID)
	})

	t.Run("", func(t *testing.T) {
		_, err = api.PostsDelete(&v, &PostsDeleteReq{PostID: postID})
		require.NoError(t, err)
	})

	t.Run("", func(t *testing.T) {
		_, err := api.PostsListByID(context.Background(), &v, &PostsListByIdReq{ID: postID})
		requireAPIError(t, err, "PostNotFound")
	})
}

func TestAPI_PostsLikes(t *testing.T) {
	api := API{}

	closer := createTestDB(t)
	defer closer()

	user := store.User{}
	_ = store2.GlobalStore.Users.Add(&user)
	v := Viewer{UserID: user.ID}

	addResp, _ := api.PostsAdd(context.Background(), &v, &PostsAddReq{Text: "test text"})

	t.Run("like post", func(t *testing.T) {
		_, err := api.PostsLike(&v, &PostsLikeReq{
			PostID: addResp.ID,
			Action: "LIKE",
		})
		require.NoError(t, err)
	})

	t.Run("like post again", func(t *testing.T) {
		_, err := api.PostsLike(&v, &PostsLikeReq{
			PostID: addResp.ID,
			Action: "LIKE",
		})
		require.NoError(t, err)
	})

	t.Run("check count and flag", func(t *testing.T) {
		listResp, err := api.PostsListByID(context.Background(), &v, &PostsListByIdReq{ID: addResp.ID})
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
		listResp, err := api.PostsListByID(context.Background(), &v, &PostsListByIdReq{ID: addResp.ID})
		require.NoError(t, err)
		require.Equal(t, 0, listResp.LikesCount)
		require.False(t, listResp.IsLiked)
	})
}

func TestAPI_PostsAdd(t *testing.T) {
	api := API{}
	closer := createTestDB(t)
	defer closer()

	_, err := api.PostsAdd(context.Background(), &Viewer{UserID: 14}, &PostsAddReq{Text: ""})
	requireAPIError(t, err, "TextEmpty")

	_, err = api.PostsAdd(context.Background(), &Viewer{UserID: 14}, &PostsAddReq{Text: strings.Repeat("a", 10000)})
	requireAPIError(t, err, "TextTooLong")

	_, err = api.PostsAdd(context.Background(), &Viewer{}, &PostsAddReq{Text: "test"})
	requireAPIError(t, err, "NotAuthorized")
}

func TestAPI_PostsListByUser(t *testing.T) {
	api := API{}
	closer := createTestDB(t)
	defer closer()

	for i := 0; i < 12; i++ {
		post := store.Post{UserID: 10}
		err := store2.GlobalStore.Posts.Add(&post)
		require.NoError(t, err)

		err = store2.GlobalStore.Wall.Add(10, post.ID)
		require.NoError(t, err)
	}

	resp, err := api.PostsList(context.Background(), &Viewer{}, &PostsListReq{
		ByUserID: "10",
		Count:    10,
	})
	require.NoError(t, err)
	require.Len(t, resp.Items, 10)
	require.NotEmpty(t, resp.PageToken)

	resp, err = api.PostsList(context.Background(), &Viewer{}, &PostsListReq{
		ByUserID:  "10",
		Count:     10,
		PageToken: resp.PageToken,
	})
	require.NoError(t, err)
	require.Len(t, resp.Items, 2)
	require.Empty(t, resp.PageToken)
}
