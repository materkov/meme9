package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
)

type Post struct {
	ID     int
	Date   int
	Text   string
	UserID int

	IsDeleted bool
}

type User struct {
	ID   int
	Name string
	Bio  string

	VkID          int
	VkAccessToken string
	VkPhoto200    string

	Email        string
	PasswordHash string
}

var RedisClient *redis.Client

type AuthToken struct {
	ID     int
	UserID int
	Token  string
	Date   int
}

type Config struct {
	VKAppID     int
	VKAppSecret string

	TelegramToken string
}

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
