package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/contextKeys"
	"log"
)

type PostStore struct {
	Cache map[int]*Post
}

func PostStoreFromCtx(ctx context.Context) *PostStore {
	return ctx.Value(contextKeys.PostStore).(*PostStore)
}

func WithPostStore(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeys.PostStore, &PostStore{Cache: map[int]*Post{}})
}

func (p *PostStore) Preload(ids []int) {
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
			post := Post{}
			err = json.Unmarshal([]byte(result.(string)), &post)
			if err != nil {
				log.Printf("Error unmarshaling post: %s", err)
			}

			p.Cache[ids[i]] = &post
		}
	}
}

func (p *PostStore) Get(id int) *Post {
	p.Preload([]int{id})
	return p.Cache[id]
}
