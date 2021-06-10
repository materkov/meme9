package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleIndex(t *testing.T) {
	setupDB(t)

	require.NoError(t, store.Followers.Add(&Followers{
		ID:         1,
		User1ID:    10,
		User2ID:    11,
		FollowDate: 1,
	}))

	require.NoError(t, store.Post.Add(&Post{ID: 1, UserID: 10}))
	require.NoError(t, store.Post.Add(&Post{ID: 2, UserID: 11}))
	require.NoError(t, store.Post.Add(&Post{ID: 3, UserID: 12}))

	resp, err := handleIndex("", &Viewer{UserID: 10})
	require.NoError(t, err)
	require.NotNil(t, resp)

	posts := resp.GetFeedRenderer().Posts
	require.Len(t, posts, 2)
	require.Equal(t, "2", posts[0].Id)
	require.Equal(t, "1", posts[1].Id)
}

func TestHandleIndex_NotAuthorized(t *testing.T) {
	setupDB(t)

	resp, err := handleIndex("", &Viewer{})
	require.NoError(t, err)
	require.Len(t, resp.GetFeedRenderer().Posts, 0)
	require.NotEmpty(t, resp.GetFeedRenderer().PlaceholderText)
}
