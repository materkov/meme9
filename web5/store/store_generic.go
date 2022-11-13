package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type CachedObject interface {
	Post | User | Photo
}

type GenericCachedStore[T CachedObject] struct {
	cache map[int]*T
}

func (p *GenericCachedStore[T]) Preload(ids []int) {
	var keys []string
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := p.cache[id]; ok {
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
			p.cache[ids[i]] = nil
		} else {
			post := new(T)
			err = json.Unmarshal([]byte(result.(string)), post)
			if err != nil {
				log.Printf("Error unmarshaling post: %s", err)
			}

			p.cache[ids[i]] = post
		}
	}
}

func (p *GenericCachedStore[T]) Get(id int) *T {
	p.Preload([]int{id})
	return p.cache[id]
}
