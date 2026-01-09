package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/materkov/meme9/likes-service/adapters/likes"
	"github.com/materkov/meme9/likes-service/api"
	likesserviceapi "github.com/materkov/meme9/likes-service/api/likes"
	likesapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/likes"
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
	likesAdapter := likes.New(client, databaseName)

	// Ensure indexes
	err = likesAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure likes indexes: %v", err)
	}

	// Initialize services
	likesService := likesserviceapi.NewService(likesAdapter)

	// Create Twirp server
	likesHandler := likesapi.NewLikesServer(likesService)
	likesHandlerWithCORS := api.AuthMiddleware(api.CORSMiddleware(likesHandler))

	// Register handler
	http.Handle(likesHandler.PathPrefix(), likesHandlerWithCORS)

	// Start HTTP server
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8084"
	}
	log.Printf("Starting likes service HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

