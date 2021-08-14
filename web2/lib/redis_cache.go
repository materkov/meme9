package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

func TryGetCached(redisClient *redis.Client, db *sqlx.DB, objType string, objId int, obj interface{}) error {
	cacheKey := fmt.Sprintf("%s:%d", objType, objId)
	cached, err := redisClient.Get(cacheKey).Bytes()
	if err == redis.Nil {
		cached = nil
	} else if err != nil {
		return fmt.Errorf("error getting from cache: %w", err)
	}

	if cached != nil {
		err = json.Unmarshal(cached, obj)
		if err != nil {
			return fmt.Errorf("error unmarshalong from cache: %w", err)
		}
		log.Printf("Get %s, found", cacheKey)
		return nil
	}

	log.Printf("Cache key %s, not found", cacheKey)

	err = db.Get(obj, fmt.Sprintf("select * from `%s` where id = %d", objType, objId))
	if err != nil {
		return fmt.Errorf("error selecting object: %w", err)
	}

	userBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error marshaling object: %w", err)
	}

	_, err = redisClient.Set(cacheKey, userBytes, time.Minute).Result()
	if err != nil {
		return fmt.Errorf("error setting cache: %w", err)
	}

	return nil
}
