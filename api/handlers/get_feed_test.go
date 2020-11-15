package handlers

import (
	"testing"

	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/store"
	"github.com/stretchr/testify/require"
)

func TestGetFeed_Handle(t *testing.T) {
	testStore, closer := InitTestStore(t)
	defer closer()

	require.NoError(t, testStore.AddPost(&store.Post{ID: 10}))
	require.NoError(t, testStore.AddPost(&store.Post{ID: 11}))

	handler := GetFeed{Store: testStore}

	// Empty feed
	resp, err := handler.Handle(&api.Viewer{}, &pb.GetFeedRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Posts, 0)

	require.NoError(t, testStore.AddToFeed(11))
	require.NoError(t, testStore.AddToFeed(10))

	// 2 posts
	resp, err = handler.Handle(&api.Viewer{}, &pb.GetFeedRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Posts, 2)
	require.Equal(t, "10", resp.Posts[0].Id)
	require.Equal(t, "11", resp.Posts[1].Id)

	// Incorrect post
	require.NoError(t, testStore.AddToFeed(1012312))

	resp, err = handler.Handle(&api.Viewer{}, &pb.GetFeedRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Posts, 2)
}
