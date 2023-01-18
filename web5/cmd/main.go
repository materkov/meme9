package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/api"
	"github.com/materkov/meme9/web5/imgproxy"
	"github.com/materkov/meme9/web5/pkg/users"
	"github.com/materkov/meme9/web5/store"
	"github.com/materkov/meme9/web5/upload"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	queue := ""
	flag.StringVar(&queue, "queue", "", "Queue listen to")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	store.RedisClient = redis.NewClient(&redis.Options{})

	configStr, err := store.RedisClient.Get(context.Background(), "config").Bytes()
	if err != nil {
		log.Fatalf("Failed reading config: %s", err)
	}

	err = json.Unmarshal(configStr, &store.DefaultConfig)
	if err != nil {
		log.Fatalf("Error parsing config JSON: %s", err)
	}

	if queue != "" {
		HandleWorker(queue)
		return
	}

	http.HandleFunc("/api2", api.HandleAPI2)
	http.HandleFunc("/upload", upload.HandleUpload)
	http.HandleFunc("/imgproxy", imgproxy.ServeHTTP)

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func HandleWorker(queue string) {
	for {
		result, err := store.RedisClient.BLPop(context.Background(), time.Second*5, queue).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return
		}

		log.Printf("Got queue task: %v", result)

		userID, _ := strconv.Atoi(result[1])
		err = users.RefreshFromVk(userID)
		if err != nil {
			log.Printf("Error doing queue: %s", err)
		}
	}
}
