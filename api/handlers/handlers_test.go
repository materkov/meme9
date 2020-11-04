package handlers

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/materkov/meme9/api/store"
	"github.com/stretchr/testify/require"
)

func InitTestStore(t *testing.T) (*store.Store, func()) {
	miniRedis, err := miniredis.Run()
	require.NoError(t, err)

	testStore := store.NewStore(redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	}))

	closer := func() {
		miniRedis.Close()
	}

	return testStore, closer
}
