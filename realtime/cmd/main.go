package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/smira/go-statsd"
	"log"
	"net/http"
	"strconv"
	"time"
)

var redisClient *redis.Client

var (
	statsdClient = statsd.NewClient("127.0.0.1:8125")

	activeConnections int
)

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

	activeConnections++
	statsdClient.Gauge("realtime_active_connections", int64(activeConnections))
	statsdClient.Incr("realtime_total_connects", 1)

	key := r.URL.Query().Get("key")
	userID, _ := strconv.Atoi(key)

	s := redisClient.Subscribe(context.Background(), fmt.Sprintf("queue:%d", userID))
	redisCh := s.Channel()

	for {
		clientDisconnect := false
		select {
		case <-r.Context().Done():
			clientDisconnect = true

		case msg, cls := <-redisCh:
			log.Printf("cls %v", cls)
			fmt.Fprintf(w, "data: %s\n", msg.Payload)
			fmt.Fprint(w, "\n")
			flusher.Flush()

		case <-time.After(time.Second * 30):
			fmt.Fprint(w, "data: {\"type\":\"ping\"}\n")
			fmt.Fprint(w, "\n")
			flusher.Flush()
		}

		if clientDisconnect {
			log.Printf("Client disconnected")
			break
		}
	}
	err := s.Close()
	if err != nil {
		log.Printf("Error closing redis pubsub: %s", err)
	}
	activeConnections--
	statsdClient.Gauge("realtime_active_connections", int64(activeConnections))
}

func Push(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("userId"))
	data := r.URL.Query().Get("data")

	_, err := redisClient.Publish(context.Background(), fmt.Sprintf("queue:%d", userID), data).Result()
	if err != nil {
		log.Printf("Error adding to redis: %s", err)
	}

	statsdClient.Incr("realtime_messages", 1)
}
