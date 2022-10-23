package testutils

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"testing"
)

func SetupRedis(t *testing.T) {
	s := miniredis.RunT(t)
	store.RedisClient = redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}