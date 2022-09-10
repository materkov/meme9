package store

import "github.com/go-redis/redis/v9"

type Post struct {
	ID     int
	Date   int
	Text   string
	UserID int
}

type User struct {
	ID   int
	Name string
	VkID int
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
}

var DefaultConfig = Config{}
