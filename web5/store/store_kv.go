package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/contextKeys"
	"log"
)

type CachedStore struct {
	Post  GenericCachedStore[Post]
	User  GenericCachedStore[User]
	Photo GenericCachedStore[Photo]
}

func CachedStoreFromCtx(ctx context.Context) *CachedStore {
	return ctx.Value(contextKeys.CachedStore).(*CachedStore)
}

func WithCachedStore(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeys.CachedStore, &CachedStore{
		Post: GenericCachedStore[Post]{
			Cache: map[int]*Post{},
		},
		User: GenericCachedStore[User]{
			Cache: map[int]*User{},
		},
		Photo: GenericCachedStore[Photo]{
			Cache: map[int]*Photo{},
		},
	})
}

type CachedObject interface {
	Post | User | Photo
}

type GenericCachedStore[T CachedObject] struct {
	Cache map[int]*T
}

func (p *GenericCachedStore[T]) Preload(ids []int) {
	var keys []string
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := p.Cache[id]; ok {
			continue
		}

		keys = append(keys, fmt.Sprintf("node:%d", id))
	}
	if len(keys) == 0 {
		return
	}

	results, err := RedisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		log.Printf("Error loading posts from redis: %s", err)
		results = make([]interface{}, len(ids))
	}

	for i, result := range results {
		if result == nil {
			p.Cache[ids[i]] = nil
		} else {
			post := new(T)
			err = json.Unmarshal([]byte(result.(string)), post)
			if err != nil {
				log.Printf("Error unmarshaling post: %s", err)
			}

			p.Cache[ids[i]] = post
		}
	}
}

func (p *GenericCachedStore[T]) Get(id int) *T {
	p.Preload([]int{id})
	return p.Cache[id]
}
