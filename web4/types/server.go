package types

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"net/http"
	"strings"
)

var RedisClient *redis.Client

func getObject(id int, obj interface{}) error {
	objBytes, err := RedisClient.Get(context.Background(), fmt.Sprintf("node:%d", id)).Bytes()
	if err == redis.Nil {
		return err
	} else if err != nil {
		return err
	}

	return json.Unmarshal(objBytes, obj)
}

func saveObject(id int, obj interface{}) error {
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

func DoHandle() {
	// CRUD words: insert, delete, update, list

	RedisClient = redis.NewClient(&redis.Options{})

	http.HandleFunc("/browse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")

		if r.Method == "OPTIONS" {
			return
		}

		authToken := r.Header.Get("authorization")
		authToken = strings.TrimPrefix(authToken, "Bearer ")
		userID, _ := AuthCheckToken(authToken)

		viewer := Viewer{
			UserID: userID,
			Origin: r.Header.Get("origin"),
		}

		resp := Browse(r.URL.Query().Get("url"), r.URL.Query().Get("q"), &viewer)
		_ = json.NewEncoder(w).Encode(resp)
	})

	http.ListenAndServe("127.0.0.1:8000", nil)
}
