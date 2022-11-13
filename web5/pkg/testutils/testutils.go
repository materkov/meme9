package testutils

import (
	"context"
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

func PrepareContext() context.Context {
	ctx := context.Background()
	ctx = store.WithCachedStore(ctx)

	return ctx
}
