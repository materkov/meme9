package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/materkov/meme9/api/pkg"
	"github.com/materkov/meme9/api/server"
	"github.com/materkov/meme9/api/store"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	conf := &pkg.Config{}

	configJson := []byte(os.Getenv("CONFIG"))
	if len(configJson) == 0 {
		homeDir, _ := os.UserHomeDir()
		configJson, _ = ioutil.ReadFile(homeDir + "/.meme")
	}

	err := json.Unmarshal(configJson, conf)
	if err != nil {
		panic("Error parsing config: " + err.Error())
	}

	redisClient := redis.NewClient(&redis.Options{})
	dataStore := store.NewStore(redisClient)

	m := server.Main{
		Store:  dataStore,
		Config: conf,
	}
	m.Run()
}
