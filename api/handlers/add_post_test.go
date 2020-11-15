package handlers

import (
	"strconv"
	"testing"

	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/store"
	"github.com/stretchr/testify/require"
)

func TestAddPost_Handle(t *testing.T) {
	testStore, closer := InitTestStore(t)
	defer closer()

	user := store.User{ID: 15}
	require.NoError(t, testStore.AddUser(&user))

	handler := AddPost{Store: testStore}

	resp, err := handler.Handle(&api.Viewer{User: &user, CSRFValidated: true}, &pb.AddPostRequest{
		Text: " test post ",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Id)
	require.Equal(t, "test post", resp.Text)

	postID, _ := strconv.Atoi(resp.Id)
	require.NotEmpty(t, postID)

	post, err := testStore.GetPost(postID)
	require.NoError(t, err)
	require.Equal(t, 15, post.UserID)

	postIds, err := testStore.GetFeed()
	require.NoError(t, err)
	require.Equal(t, []int{postID}, postIds)
}
