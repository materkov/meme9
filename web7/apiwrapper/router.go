package apiwrapper

import (
	"log"
	"net/http"

	"github.com/materkov/meme9/web7/api"
)

type Router struct {
	baseHandler      *BaseHandler
	feedHandler      *FeedHandler
	publishHandler   *PublishHandler
	loginHandler     *LoginHandler
	registerHandler  *RegisterHandler
	userPostsHandler *UserPostsHandler
	subscribeHandler *SubscribeHandler
}

func NewRouter(api *api.API) *Router {
	return &Router{
		baseHandler:      NewBaseHandler(api),
		feedHandler:      NewFeedHandler(api),
		publishHandler:   NewPublishHandler(api),
		loginHandler:     NewLoginHandler(api),
		registerHandler:  NewRegisterHandler(api),
		userPostsHandler: NewUserPostsHandler(api),
		subscribeHandler: NewSubscribeHandler(api),
	}
}

func (r *Router) RegisterRoutes() {

	// API Endpoints (JSON responses)
	http.HandleFunc("/api/feed", CORSMiddleware(JSONMiddleware(r.feedHandler.Handle)))
	http.HandleFunc("/api/publish", CORSMiddleware(JSONMiddleware(r.baseHandler.AuthMiddleware(r.publishHandler.Handle))))
	http.HandleFunc("/api/login", CORSMiddleware(JSONMiddleware(r.loginHandler.Handle)))
	http.HandleFunc("/api/register", CORSMiddleware(JSONMiddleware(r.registerHandler.Handle)))
	http.HandleFunc("/api/userPosts", CORSMiddleware(JSONMiddleware(r.userPostsHandler.Handle)))
	http.HandleFunc("/api/subscribe", CORSMiddleware(JSONMiddleware(r.baseHandler.AuthMiddleware(r.subscribeHandler.HandleSubscribe))))
	http.HandleFunc("/api/unsubscribe", CORSMiddleware(JSONMiddleware(r.baseHandler.AuthMiddleware(r.subscribeHandler.HandleUnsubscribe))))
	http.HandleFunc("/api/subscriptionStatus", CORSMiddleware(JSONMiddleware(r.baseHandler.AuthMiddleware(r.subscribeHandler.HandleSubscriptionStatus))))
}

func (r *Router) StartServer(addr string) {
	log.Printf("Starting HTTP server at http://%s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
