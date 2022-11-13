package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
)

var RedisClient *redis.Client

var DefaultConfig = Config{}

var ErrNodeNotFound = fmt.Errorf("node not found")

func NodeGet(id int, obj interface{}) error {
	objBytes, err := RedisClient.Get(context.Background(), fmt.Sprintf("node:%d", id)).Bytes()
	if err == redis.Nil {
		return ErrNodeNotFound
	} else if err != nil {
		return fmt.Errorf("error getting node from redis: %s", err)
	}

	return json.Unmarshal(objBytes, obj)
}

func NodeSave(id int, obj interface{}) error {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = RedisClient.Set(context.Background(), fmt.Sprintf("node:%d", id), objBytes, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
