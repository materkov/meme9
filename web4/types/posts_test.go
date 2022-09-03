package types

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func setupRedis(t *testing.T) {
	s := miniredis.RunT(t)
	redisClient = redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}

func TestPosts_CRUD(t *testing.T) {
	setupRedis(t)

	// Add post
	req := &postsAddRequest{
		Text: "test post",
	}
	postID, err := postsAdd(req, &Viewer{UserID: 161})
	require.NoError(t, err)
	require.NotEmpty(t, postID)

	// List posts
	respList := postsList([]int{postID, postID + 1})
	require.Len(t, respList, 2)

	require.Equal(t, "test post", respList[0].Text)
	require.Equal(t, strconv.Itoa(postID), respList[0].ID)
	require.Equal(t, "161", respList[0].FromID)
	require.Equal(t, fmt.Sprintf("/posts/%d", postID), respList[0].DetailsURL)
	require.Equal(t, strconv.Itoa(postID+1), respList[1].ID)
}

func TestPosts_AddErrors(t *testing.T) {
	setupRedis(t)

	_, err := postsAdd(&postsAddRequest{}, &Viewer{UserID: 0})
	require.ErrorContains(t, err, "zero viewer")
}

func TestPosts_Feed(t *testing.T) {
	setupRedis(t)

	post1, err1 := postsAdd(&postsAddRequest{}, &Viewer{UserID: 15})
	post2, err2 := postsAdd(&postsAddRequest{}, &Viewer{UserID: 16})
	post3, err3 := postsAdd(&postsAddRequest{}, &Viewer{UserID: 16})

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	require.NotEmpty(t, post1)
	require.NotEmpty(t, post2)
	require.NotEmpty(t, post3)

	postIds, err := postsGetFeed()
	require.NoError(t, err)
	require.Equal(t, []int{post3, post2, post1}, postIds)

	postIds, err = postsGetFeedByUsers([]int{16, 17})
	require.NoError(t, err)
	require.Equal(t, []int{post3, post2}, postIds)
}
