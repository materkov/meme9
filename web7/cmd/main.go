package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/materkov/meme9/web7/adapters/likes"
	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	"github.com/materkov/meme9/web7/api/auth"
	likesserviceapi "github.com/materkov/meme9/web7/api/likes"
	postsserviceapi "github.com/materkov/meme9/web7/api/posts"
	subscriptionsserviceapi "github.com/materkov/meme9/web7/api/subscriptions"
	usersserviceapi "github.com/materkov/meme9/web7/api/users"
	authapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/auth"
	likesapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/likes"
	postsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/posts"
	subscriptionsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/subscriptions"
	usersapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/users"
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
	databaseName := "meme9"
	postsAdapter := posts.New(client, databaseName)
	usersAdapter := users.New(client, databaseName)
	tokensAdapter := tokens.New(client, databaseName)
	subscriptionsAdapter := subscriptions.New(client, databaseName)
	likesAdapter := likes.New(client, databaseName)

	// Ensure indexes
	err = usersAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure user indexes: %v", err)
	}
	err = subscriptionsAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure subscription indexes: %v", err)
	}
	err = likesAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure likes indexes: %v", err)
	}

	// Initialize services
	tokensService := tokensservice.New(tokensAdapter)

	// Create separate service instances
	likesService := likesserviceapi.NewService(likesAdapter, usersAdapter)
	postsServiceInstance := postsserviceapi.NewService(postsAdapter, usersAdapter, subscriptionsAdapter, likesAdapter)
	authService := auth.NewService(usersAdapter, tokensAdapter, tokensService)
	usersService := usersserviceapi.NewService(usersAdapter)
	subscriptionsService := subscriptionsserviceapi.NewService(subscriptionsAdapter)

	// Create Twirp servers for each service

	postsHandler := postsapi.NewPostsServer(postsServiceInstance)
	// Auth service should NOT have authHooks applied - it handles its own validation
	// and VerifyToken is called from within the hook, which would cause infinite recursion
	authHandler := authapi.NewAuthServer(authService)
	usersHandler := usersapi.NewUsersServer(usersService)
	subscriptionsHandler := subscriptionsapi.NewSubscriptionsServer(subscriptionsService)
	likesHandler := likesapi.NewLikesServer(likesService)

	// Wrap with CORS middleware
	postsHandlerWithCORS := api.AuthMiddleware(authService, api.CORSMiddleware(postsHandler))
	authHandlerWithCORS := api.AuthMiddleware(authService, api.CORSMiddleware(authHandler))
	usersHandlerWithCORS := api.AuthMiddleware(authService, api.CORSMiddleware(usersHandler))
	subscriptionsHandlerWithCORS := api.AuthMiddleware(authService, api.CORSMiddleware(subscriptionsHandler))
	likesHandlerWithCORS := api.AuthMiddleware(authService, api.CORSMiddleware(likesHandler))

	// Register all handlers
	http.Handle(postsHandler.PathPrefix(), postsHandlerWithCORS)
	http.Handle(authHandler.PathPrefix(), authHandlerWithCORS)
	http.Handle(usersHandler.PathPrefix(), usersHandlerWithCORS)
	http.Handle(subscriptionsHandler.PathPrefix(), subscriptionsHandlerWithCORS)
	http.Handle(likesHandler.PathPrefix(), likesHandlerWithCORS)

	// Start HTTP server
	addr := "127.0.0.1:8080"
	log.Printf("Starting HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
