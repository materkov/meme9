package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/materkov/meme9/web7/adapters/mongo"
)

type API struct {
	mongo *mongo.Adapter
}

func NewAPI(mongo *mongo.Adapter) *API {
	return &API{mongo: mongo}
}

type Post struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAd"`
}

func mapMongoPostToAPIPost(mongoPost mongo.Post) Post {
	return Post{
		ID:        mongoPost.ID,
		Text:      mongoPost.Text,
		CreatedAt: mongoPost.CreatedAt.Format(time.RFC3339),
	}
}

func (a *API) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (a *API) feedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	mongoPosts, err := a.mongo.GetAllPosts(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}

	apiPosts := make([]Post, len(mongoPosts))
	for i, mongoPost := range mongoPosts {
		apiPosts[i] = mapMongoPostToAPIPost(mongoPost)
	}

	json.NewEncoder(w).Encode(apiPosts)
}

type PublishReq struct {
	Text string `json:"text"`
}

type PublishResp struct {
	ID string `json:"id"`
}

func (a *API) publishHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	var publishReq PublishReq
	err = json.Unmarshal(body, &publishReq)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	post, err := a.mongo.AddPost(r.Context(), mongo.Post{
		Text:      publishReq.Text,
		CreatedAt: time.Now(),
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create post"})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(PublishResp{ID: post.ID})
}

func (a *API) Serve() {
	http.HandleFunc("/feed", a.corsMiddleware(a.feedHandler))
	http.HandleFunc("/publish", a.corsMiddleware(a.publishHandler))

	log.Printf("Starting HTTP server at http://127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
