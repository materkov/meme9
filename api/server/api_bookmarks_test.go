package server

import (
	"context"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestBookmarks_CRUD(t *testing.T) {
	_ = createTestDB(t)
	viewer := Viewer{UserID: 16}
	ctx := context.WithValue(context.Background(), CtxViewerKey, &viewer)

	post := store.Post{}
	_ = store2.GlobalStore.Posts.Add(&post)

	// Add
	apiSrv := BookmarkServer{}
	respAdd, err := apiSrv.Add(ctx, &api.BookmarksAddReq{PostId: strconv.Itoa(post.ID)})
	require.NoError(t, err)
	require.NotNil(t, respAdd)

	// Check
	respList, err := apiSrv.List(ctx, &api.BookmarkListReq{})
	require.NoError(t, err)
	require.Empty(t, respList.PageToken)
	require.Len(t, respList.Items, 1)
	require.Equal(t, respList.Items[0].Post.Id, strconv.Itoa(post.ID))
	require.True(t, respList.Items[0].Post.IsBookmarked)

	// Remove
	respRemove, err := apiSrv.Remove(ctx, &api.BookmarksAddReq{PostId: strconv.Itoa(post.ID)})
	require.NoError(t, err)
	require.NotNil(t, respRemove)

	// Check again
	respList, err = apiSrv.List(ctx, &api.BookmarkListReq{})
	require.NoError(t, err)
	require.Empty(t, respList.PageToken)
	require.Len(t, respList.Items, 0)
}

func TestBookmarks_DeletedPost(t *testing.T) {
	_ = createTestDB(t)
	ctx := createViewerContext(199)
	srv := BookmarkServer{}

	// Create post
	post := store.Post{Text: "test post"}
	_ = store2.GlobalStore.Posts.Add(&post)

	// Add bookmark
	_, err := srv.Add(ctx, &api.BookmarksAddReq{PostId: strconv.Itoa(post.ID)})
	require.NoError(t, err)

	// Delete post
	post.IsDeleted = true
	_ = store2.GlobalStore.Posts.Update(&post)

	// Check bookmarks list
	listResp, err := srv.List(ctx, &api.BookmarkListReq{})
	require.NoError(t, err)
	require.True(t, listResp.Items[0].Post.IsDeleted)
	require.Empty(t, listResp.Items[0].Post.Text)

	// Try to add a deleted post
	_, err = srv.Add(ctx, &api.BookmarksAddReq{PostId: strconv.Itoa(post.ID)})
	requireAPIError(t, err, "PostNotFound")
}
