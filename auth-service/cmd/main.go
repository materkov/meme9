package main

import (
	"context"
	"log"
	"net/http"
	"os"

	authapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/auth"
	"github.com/materkov/meme9/auth-service/adapters/tokens"
	"github.com/materkov/meme9/auth-service/adapters/users"
	"github.com/materkov/meme9/auth-service/api"
	tokensservice "github.com/materkov/meme9/auth-service/services/tokens"
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
	usersAdapter := users.New(client, databaseName)
	tokensAdapter := tokens.New(client, databaseName)

	// Ensure indexes
	err = usersAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure user indexes: %v", err)
	}

	// Initialize services
	tokensService := tokensservice.New(tokensAdapter)
	authService := api.NewService(usersAdapter, tokensAdapter, tokensService)

	// Create Twirp server
	authHandler := authapi.NewAuthServer(authService)
	authHandlerWithCORS := api.AuthMiddleware(authService, api.CORSMiddleware(authHandler))

	// Register handler
	http.Handle(authHandler.PathPrefix(), authHandlerWithCORS)

	// Start HTTP server
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	log.Printf("Starting auth service HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
