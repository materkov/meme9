package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"time"
)

var redisClient *redis.Client
var msgId int

func main() {
	redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Cannot connect to redis: %s", err)
	}

	log.Print("Starting HTTP server at 127.0.0.1:8001")

	http.HandleFunc("/listen", Server)
	http.HandleFunc("/push", Push)

	_ = http.ListenAndServe("127.0.0.1:8001", nil)
}

func Server(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	flusher.Flush()

	key := r.URL.Query().Get("key")
	userID, _ := strconv.Atoi(key)

	for {
		items, err := redisClient.BLPop(context.Background(), time.Second*5, fmt.Sprintf("queue:%d", userID)).Result()
		if errors.Is(err, redis.Nil) {
			log.Printf("Got nil from redis")
			// No data
		} else if err != nil {
			log.Printf("Error getting from redis: %s", err)
		} else {
			msgId++
			fmt.Fprintf(w, "id: %d\n", msgId)
			fmt.Fprintf(w, "data: %s\n", items[1])
			fmt.Fprintf(w, "\n")
			flusher.Flush()
		}

		isClosed := false
		select {
		case <-r.Context().Done():
			isClosed = true
		default:
		}

		if isClosed {
			log.Printf("Client closed connection")
			break
		}

		log.Printf("Client alive")
	}
}

func Push(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))
	data := r.URL.Query().Get("data")

	_, err := redisClient.RPush(context.Background(), fmt.Sprintf("queue:%d", userID), data).Result()
	if err != nil {
		log.Printf("Error adding to redis: %s", err)
	}
}
