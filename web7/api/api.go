package api

import (
	"log"
	"net/http"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	postsservice "github.com/materkov/meme9/web7/services/posts"
	tokensservice "github.com/materkov/meme9/web7/services/tokens"
)

type API struct {
	posts         *posts.Adapter
	users         *users.Adapter
	tokens        *tokens.Adapter
	subscriptions *subscriptions.Adapter

	postsService  *postsservice.Service
	tokensService *tokensservice.Service
}

func NewAPI(postsAdapter *posts.Adapter, usersAdapter *users.Adapter, tokensAdapter *tokens.Adapter, subscriptionsAdapter *subscriptions.Adapter, postsService *postsservice.Service, tokensService *tokensservice.Service) *API {
	return &API{
		posts:         postsAdapter,
		users:         usersAdapter,
		tokens:        tokensAdapter,
		subscriptions: subscriptionsAdapter,
		postsService:  postsService,
		tokensService: tokensService,
	}
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

func (a *API) Serve() {
	http.HandleFunc("/api/feed", corsMiddleware(jsonMiddleware(a.feedHandler)))
	http.HandleFunc("/api/publish", corsMiddleware(jsonMiddleware(a.authMiddleware(a.publishHandler))))
	http.HandleFunc("/api/login", corsMiddleware(jsonMiddleware(a.loginHandler)))
	http.HandleFunc("/api/register", corsMiddleware(jsonMiddleware(a.registerHandler)))
	http.HandleFunc("/api/userPosts", corsMiddleware(jsonMiddleware(a.userPostsHandler)))
	http.HandleFunc("/api/subscribe", corsMiddleware(jsonMiddleware(a.authMiddleware(a.subscribeHandler))))
	http.HandleFunc("/api/unsubscribe", corsMiddleware(jsonMiddleware(a.authMiddleware(a.unsubscribeHandler))))
	http.HandleFunc("/api/subscriptionStatus", corsMiddleware(jsonMiddleware(a.authMiddleware(a.subscriptionStatusHandler))))
	http.HandleFunc("/users/{id}", a.userPageHandler)
	http.HandleFunc("/posts/{id}", a.postPageHandler)
	http.HandleFunc("/feed", a.feedPageHandler)
	http.HandleFunc("/favicon.ico", a.faviconHandler)
	http.HandleFunc("/static/", a.staticHandler)
	http.HandleFunc("/", a.indexHandler)

	log.Printf("Starting HTTP server at http://127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
