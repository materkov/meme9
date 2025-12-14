package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	authapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/auth"
	feedapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/feed"
	postsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/posts"
	subscriptionsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/subscriptions"
	usersapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/users"
	postsservice "github.com/materkov/meme9/web7/services/posts"
	tokensservice "github.com/materkov/meme9/web7/services/tokens"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@localhost:27017/meme9?authSource=admin"
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping MongoDB to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")

	// Initialize adapters
	postsAdapter := posts.New(client)
	usersAdapter := users.New(client)
	tokensAdapter := tokens.New(client)
	subscriptionsAdapter := subscriptions.New(client)

	// Ensure indexes
	err = usersAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure user indexes: %v", err)
	}
	err = subscriptionsAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure subscription indexes: %v", err)
	}

	// Initialize services
	postsService := postsservice.New(postsAdapter)
	tokensService := tokensservice.New(tokensAdapter)

	// Create separate service instances
	feedService := api.NewFeedService(postsAdapter, usersAdapter, subscriptionsAdapter)
	postsServiceInstance := api.NewPostsService(postsAdapter, usersAdapter, postsService)
	authService := api.NewAuthService(usersAdapter, tokensAdapter, tokensService)
	usersService := api.NewUsersService(usersAdapter)
	subscriptionsService := api.NewSubscriptionsService(subscriptionsAdapter)

	// Create Twirp servers for each service
	authHooks := api.AuthHook(authService)

	feedHandler := feedapi.NewFeedServer(feedService, twirp.WithServerHooks(authHooks))
	postsHandler := postsapi.NewPostsServer(postsServiceInstance, twirp.WithServerHooks(authHooks))
	// Auth service should NOT have authHooks applied - it handles its own validation
	// and VerifyToken is called from within the hook, which would cause infinite recursion
	authHandler := authapi.NewAuthServer(authService)
	usersHandler := usersapi.NewUsersServer(usersService, twirp.WithServerHooks(authHooks))
	subscriptionsHandler := subscriptionsapi.NewSubscriptionsServer(subscriptionsService, twirp.WithServerHooks(authHooks))

	// Wrap with auth middleware to inject headers into context
	feedHandlerWithAuth := api.AuthMiddleware(feedHandler)
	postsHandlerWithAuth := api.AuthMiddleware(postsHandler)
	authHandlerWithAuth := api.AuthMiddleware(authHandler)
	usersHandlerWithAuth := api.AuthMiddleware(usersHandler)
	subscriptionsHandlerWithAuth := api.AuthMiddleware(subscriptionsHandler)

	// Wrap with CORS middleware
	feedHandlerWithCORS := api.CORSMiddleware(feedHandlerWithAuth)
	postsHandlerWithCORS := api.CORSMiddleware(postsHandlerWithAuth)
	authHandlerWithCORS := api.CORSMiddleware(authHandlerWithAuth)
	usersHandlerWithCORS := api.CORSMiddleware(usersHandlerWithAuth)
	subscriptionsHandlerWithCORS := api.CORSMiddleware(subscriptionsHandlerWithAuth)

	// Register all handlers
	http.Handle(feedHandler.PathPrefix(), feedHandlerWithCORS)
	http.Handle(postsHandler.PathPrefix(), postsHandlerWithCORS)
	http.Handle(authHandler.PathPrefix(), authHandlerWithCORS)
	http.Handle(usersHandler.PathPrefix(), usersHandlerWithCORS)
	http.Handle(subscriptionsHandler.PathPrefix(), subscriptionsHandlerWithCORS)

	// Start HTTP server
	addr := "127.0.0.1:8080"
	log.Printf("Starting HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
