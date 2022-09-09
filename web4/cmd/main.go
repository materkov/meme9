package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web4/api"
	"github.com/materkov/meme9/web4/types"
	"math/rand"
	"os"
	"time"
)

func main() {
	types.RedisClient = redis.NewClient(&redis.Options{})
	
	rand.Seed(time.Now().UnixNano())

	homeDir, _ := os.UserHomeDir()
	dat, _ := os.ReadFile(homeDir + "/mypage/config.json")
	if len(dat) > 0 {
		_ = json.Unmarshal(dat, &types.DefaultConfig)
	}

	config := os.Getenv("CONFIG")
	if config != "" {
		_ = json.Unmarshal([]byte(config), &types.DefaultConfig)
	}

	api.Serve()
}
