package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/contextKeys"
	"log"
)

type UserStore struct {
	Cache map[int]*User
}

func UserStoreFromCtx(ctx context.Context) *UserStore {
	return ctx.Value(contextKeys.UserStore).(*UserStore)
}

func WithUserStore(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeys.UserStore, &UserStore{
		Cache: map[int]*User{},
	})
}

func (p *UserStore) Preload(ids []int) {
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
		log.Printf("Error loading from redis: %s", err)
		results = make([]interface{}, len(ids))
	}

	for i, result := range results {
		if result == nil {
			p.Cache[ids[i]] = nil
		} else {
			obj := User{}
			err = json.Unmarshal([]byte(result.(string)), &obj)
			if err != nil {
				log.Printf("Error unmarshaling user: %s", err)
			}

			p.Cache[ids[i]] = &obj
		}
	}
}

func (p *UserStore) Get(id int) *User {
	p.Preload([]int{id})
	return p.Cache[id]
}
