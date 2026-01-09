package main

import (
	"context"
	"log"
	"net/http"
	"os"

	postsapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/posts"
	"github.com/materkov/meme9/posts-service/adapters/posts"
	"github.com/materkov/meme9/posts-service/api"
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

	// Initialize services
	postsService := api.NewService(postsAdapter)

	// Create Twirp server
	postsHandler := postsapi.NewPostsServer(postsService)
	postsHandlerWithCORS := api.AuthMiddleware(api.CORSMiddleware(postsHandler))

	// Register handler
	http.Handle(postsHandler.PathPrefix(), postsHandlerWithCORS)

	// Start HTTP server
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8085"
	}
	log.Printf("Starting posts service HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
