package api

import (
	"log"
	"net/http"

	"github.com/materkov/meme9/web7/adapters/mongo"
)

type API struct {
	mongo *mongo.Adapter
}

func NewAPI(mongo *mongo.Adapter) *API {
	return &API{mongo: mongo}
}

func jsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func (a *API) Serve() {
	http.HandleFunc("/feed", corsMiddleware(jsonMiddleware(a.feedHandler)))
	http.HandleFunc("/publish", corsMiddleware(jsonMiddleware(a.publishHandler)))
	http.HandleFunc("/login", corsMiddleware(jsonMiddleware(a.loginHandler)))
	http.HandleFunc("/register", corsMiddleware(jsonMiddleware(a.registerHandler)))
	http.HandleFunc("/static/", a.staticHandler)
	http.HandleFunc("/", a.indexHandler)

	log.Printf("Starting HTTP server at http://127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
