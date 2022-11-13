package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"strconv"
)

type likedData struct {
	isLiked bool
	count   int
}

type LikedStore struct {
	cache map[string]likedData
}

func (l *LikedStore) Preload(viewerID int, postIds []int) {
	pipe := RedisClient.Pipeline()

	scoresCmd := map[int]*redis.FloatCmd{}
	cardCmd := map[int]*redis.IntCmd{}

	for _, id := range postIds {
		if id <= 0 {
			continue
		}

		cacheKey := fmt.Sprintf("%d:%d", viewerID, id)
		if _, ok := l.cache[cacheKey]; ok {
			continue
		}

		if viewerID > 0 {
			scoresCmd[id] = pipe.ZScore(context.Background(), fmt.Sprintf("postLikes:%d", id), strconv.Itoa(viewerID))
		}
		cardCmd[id] = pipe.ZCard(context.Background(), fmt.Sprintf("postLikes:%d", id))
	}

	if len(cardCmd) == 0 {
		return
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error loading likes from redis: %s", err)
	}

	for _, id := range postIds {
		if id <= 0 {
			continue
		}

		cacheKey := fmt.Sprintf("%d:%d", viewerID, id)

		count := int(cardCmd[id].Val())

		isLiked := false
		if viewerID > 0 {
			isLiked = scoresCmd[id].Val() > 0
		}

		l.cache[cacheKey] = likedData{isLiked: isLiked, count: count}
	}
}

func (l *LikedStore) Get(viewerID int, postID int) (bool, int) {
	l.Preload(viewerID, []int{postID})

	cacheKey := fmt.Sprintf("%d:%d", viewerID, postID)
	data := l.cache[cacheKey]
	return data.isLiked, data.count
}
