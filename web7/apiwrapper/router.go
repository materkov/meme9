package apiwrapper

import (
	"log"
	"net/http"

	"github.com/materkov/meme9/web7/api"
)

type Router struct {
	api *api.API
}

func NewRouter(api *api.API) *Router {
	return &Router{api: api}
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (r *Router) RegisterRoutes() {
	// API Endpoints (JSON responses)
	http.HandleFunc("/api/feed", corsMiddleware(jsonMiddleware(r.feedHandler)))
	http.HandleFunc("/api/publish", corsMiddleware(jsonMiddleware(r.authMiddleware(r.publishHandler))))
	http.HandleFunc("/api/login", corsMiddleware(jsonMiddleware(r.loginHandler)))
	http.HandleFunc("/api/register", corsMiddleware(jsonMiddleware(r.registerHandler)))
	http.HandleFunc("/api/userPosts", corsMiddleware(jsonMiddleware(r.userPostsHandler)))
	http.HandleFunc("/api/subscribe", corsMiddleware(jsonMiddleware(r.authMiddleware(r.subscribeHandler))))
	http.HandleFunc("/api/unsubscribe", corsMiddleware(jsonMiddleware(r.authMiddleware(r.unsubscribeHandler))))
	http.HandleFunc("/api/subscriptionStatus", corsMiddleware(jsonMiddleware(r.authMiddleware(r.subscriptionStatusHandler))))
}

func (r *Router) StartServer(addr string) {
	log.Printf("Starting HTTP server at http://%s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
