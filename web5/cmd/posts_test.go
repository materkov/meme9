package main

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func setupRedis(t *testing.T) {
	s := miniredis.RunT(t)
	store.RedisClient = redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}

func TestPosts_CRUD(t *testing.T) {
	setupRedis(t)

	postID, err := postsAdd("test post", 200)
	require.NoError(t, err)
	require.NotEmpty(t, postID)

	posts := postsList([]int{postID, 0})
	require.Len(t, posts, 2)

	require.Equal(t, strconv.Itoa(postID), posts[0].ID)
	require.Equal(t, "test post", posts[0].Text)
	require.NotEmpty(t, posts[0].Date)
	require.Equal(t, "200", posts[0].UserID)

	require.Equal(t, "0", posts[1].ID)
}
