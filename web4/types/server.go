package types

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"net/http"
	"time"
)

var redisClient *redis.Client

func getObject(id int, obj interface{}) error {
	objBytes, err := redisClient.Get(context.Background(), fmt.Sprintf("node:%d", id)).Bytes()
	if err == redis.Nil {
		return err
	} else if err != nil {
		return err
	}

	return json.Unmarshal(objBytes, obj)
}

func DoHandle() {
	// CRUD words: insert, delete, update, list
	//db, err := sql.Open("mysql", "root:root@/meme9")
	//if err != nil {
	//	log.Fatalf("mysql error: %s", err)
	//}

	redisClient = redis.NewClient(&redis.Options{})

	http.HandleFunc("/browse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		time.Sleep(time.Second * 1)
		resp := Browse(r.URL.Query().Get("url"))
		_ = json.NewEncoder(w).Encode(resp)
	})
	http.HandleFunc("/posts.insert", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		req := AddPostRequest{}
		json.NewDecoder(r.Body).Decode(&req)

		resp := addPost(&req)
		_ = json.NewEncoder(w).Encode(resp)
	})

	http.ListenAndServe("127.0.0.1:8000", nil)
}
