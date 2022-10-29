package store

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/contextKeys"
	"log"
	"time"
)

type OnlineStore struct {
	cache  map[int]bool
	needed map[int]bool
}

func OnlineStoreFromCtx(ctx context.Context) *OnlineStore {
	return ctx.Value(contextKeys.OnlineStore).(*OnlineStore)
}

func WithOnlineStore(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeys.OnlineStore, &OnlineStore{
		cache:  map[int]bool{},
		needed: map[int]bool{},
	})
}

func (o *OnlineStore) Preload(id int) {
	if _, ok := o.cache[id]; ok {
		return
	}
	o.needed[id] = true
}

func (o *OnlineStore) Get(id int) bool {
	o.Preload(id)

	if len(o.needed) > 0 {
		neededKeys := make([]string, len(o.needed))
		neededIds := make([]int, len(o.needed))

		idx := 0
		for userID := range o.needed {
			neededKeys[idx] = fmt.Sprintf("online:%d", userID)
			neededIds[idx] = userID
			idx++
		}
		o.needed = nil

		results, err := RedisClient.MGet(context.Background(), neededKeys...).Result()
		if err != nil {
			log.Printf("Error getting online from redis: %s", err)
		} else {
			idx = 0
			for _, userID := range neededIds {
				o.cache[userID] = results[idx] != nil
				idx++
			}
		}

		o.needed = map[int]bool{}
	}

	return o.cache[id]
}

func (o *OnlineStore) Set(id int) error {
	_, err := RedisClient.Set(context.Background(), fmt.Sprintf("online:%d", id), time.Now().Unix(), time.Minute*3).Result()
	return err
}
