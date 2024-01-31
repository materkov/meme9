package server

import (
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
	"strconv"
	"strings"
	"testing"
)

func requireAPIError(t *testing.T, err error, code string) {
	t.Helper()
	if twirpErr, ok := err.(twirp.Error); ok {
		require.Equal(t, code, twirpErr.Msg())
	} else {
		require.Fail(t, "err is not twirp error")
	}
}

func TestAPI_PostsCRUD(t *testing.T) {
	srv := PostsServer{}
	closer := createTestDB(t)
	defer closer()

	user := store.User{}
	_ = store2.GlobalStore.Users.Add(&user)
	ctx := createViewerContext(user.ID)

	addResp, err := srv.Add(ctx, &api.AddReq{
		Text: "test text",
	})
	require.NoError(t, err)
	require.NotNil(t, addResp)
	require.NotEmpty(t, addResp.Id)

	postID := addResp.Id

	t.Run("", func(t *testing.T) {
		resp, err := srv.List(ctx, &api.ListReq{})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, resp.Items[0].Id, postID)
	})

	t.Run("", func(t *testing.T) {
		resp, err := srv.List(ctx, &api.ListReq{ById: postID})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, resp.Items[0].Id, postID)
	})

	t.Run("", func(t *testing.T) {
		resp, err := srv.List(ctx, &api.ListReq{ByUserId: strconv.Itoa(user.ID)})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, resp.Items[0].Id, postID)
	})

	t.Run("", func(t *testing.T) {
		_, err = srv.Delete(ctx, &api.PostsDeleteReq{PostId: postID})
		require.NoError(t, err)
	})

	t.Run("", func(t *testing.T) {
		resp, err := srv.List(ctx, &api.ListReq{ById: postID})
		require.NoError(t, err)
		require.True(t, resp.Items[0].IsDeleted)
	})
}

func TestAPI_PostsLikes(t *testing.T) {
	srv := PostsServer{}

	closer := createTestDB(t)
	defer closer()

	ctx := createViewerContext(2423)

	post := store.Post{}
	_ = store2.GlobalStore.Posts.Add(&post)

	t.Run("like post", func(t *testing.T) {
		_, err := srv.Like(ctx, &api.PostsLikeReq{
			PostId: strconv.Itoa(post.ID),
			Action: api.PostLikeAction_LIKE,
		})
		require.NoError(t, err)
	})

	t.Run("like post again", func(t *testing.T) {
		_, err := srv.Like(ctx, &api.PostsLikeReq{
			PostId: strconv.Itoa(post.ID),
			Action: api.PostLikeAction_LIKE,
		})
		require.NoError(t, err)
	})

	t.Run("check count and flag", func(t *testing.T) {
		listResp, err := srv.List(ctx, &api.ListReq{ById: strconv.Itoa(post.ID)})
		require.NoError(t, err)
		require.Equal(t, int32(1), listResp.Items[0].LikesCount)
		require.True(t, listResp.Items[0].IsLiked)
	})

	t.Run("dislike post", func(t *testing.T) {
		_, err := srv.Like(ctx, &api.PostsLikeReq{
			PostId: strconv.Itoa(post.ID),
			Action: api.PostLikeAction_UNLIKE,
		})
		require.NoError(t, err)
	})

	t.Run("check again", func(t *testing.T) {
		listResp, err := srv.List(ctx, &api.ListReq{ById: strconv.Itoa(post.ID)})
		require.NoError(t, err)
		require.Equal(t, int32(0), listResp.Items[0].LikesCount)
		require.False(t, listResp.Items[0].IsLiked)
	})

	t.Run("like not existed post", func(t *testing.T) {
		_, err := srv.Like(ctx, &api.PostsLikeReq{
			PostId: "9124679158",
			Action: api.PostLikeAction_LIKE,
		})
		requireAPIError(t, err, "PostNotFound")
	})

	t.Run("like deleted post", func(t *testing.T) {
		deletedPost := store.Post{IsDeleted: true}
		_ = store2.GlobalStore.Posts.Add(&deletedPost)

		_, err := srv.Like(ctx, &api.PostsLikeReq{
			PostId: strconv.Itoa(deletedPost.ID),
			Action: api.PostLikeAction_LIKE,
		})
		requireAPIError(t, err, "PostNotFound")
	})

	t.Run("dislike deleted or not existent posts", func(t *testing.T) {
		deletedPost := store.Post{IsDeleted: true}
		_ = store2.GlobalStore.Posts.Add(&deletedPost)

		_, err := srv.Like(ctx, &api.PostsLikeReq{
			PostId: strconv.Itoa(deletedPost.ID),
			Action: api.PostLikeAction_UNLIKE,
		})
		require.NoError(t, err)

		_, err = srv.Like(ctx, &api.PostsLikeReq{
			PostId: "452451352132",
			Action: api.PostLikeAction_UNLIKE,
		})
		require.NoError(t, err)
	})
}

func TestAPI_PostsAdd(t *testing.T) {
	srv := PostsServer{}
	closer := createTestDB(t)
	defer closer()

	_, err := srv.Add(createViewerContext(13), &api.AddReq{Text: ""})
	requireAPIError(t, err, "TextEmpty")

	_, err = srv.Add(createViewerContext(13), &api.AddReq{Text: strings.Repeat("a", 10000)})
	requireAPIError(t, err, "TextTooLong")

	_, err = srv.Add(createViewerContext(0), &api.AddReq{Text: "test"})
	requireAPIError(t, err, "NotAuthorized")
}

func TestAPI_PostsListByUser(t *testing.T) {
	srv := PostsServer{}
	closer := createTestDB(t)
	defer closer()

	for i := 0; i < 12; i++ {
		post := store.Post{UserID: 10}
		err := store2.GlobalStore.Posts.Add(&post)
		require.NoError(t, err)

		err = store2.GlobalStore.Wall.Add(10, post.ID)
		require.NoError(t, err)
	}

	resp, err := srv.List(createViewerContext(0), &api.ListReq{
		ByUserId: "10",
		Count:    10,
	})
	require.NoError(t, err)
	require.Len(t, resp.Items, 10)
	require.NotEmpty(t, resp.PageToken)

	resp, err = srv.List(createViewerContext(0), &api.ListReq{
		ByUserId:  "10",
		Count:     10,
		PageToken: resp.PageToken,
	})
	require.NoError(t, err)
	require.Len(t, resp.Items, 2)
	require.Empty(t, resp.PageToken)
}
