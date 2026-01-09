package main

import (
	"context"
	"log"
	"net/http"
	"os"

	subscriptionsapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/subscriptions"
	"github.com/materkov/meme9/subscriptions-service/adapters/subscriptions"
	"github.com/materkov/meme9/subscriptions-service/api"
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
	subscriptionsAdapter := subscriptions.New(client, databaseName)

	// Ensure indexes
	err = subscriptionsAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure subscription indexes: %v", err)
	}

	// Initialize services
	subscriptionsService := api.NewService(subscriptionsAdapter)

	// Create Twirp server
	subscriptionsHandler := subscriptionsapi.NewSubscriptionsServer(subscriptionsService)
	subscriptionsHandlerWithAuth := api.AuthMiddleware(subscriptionsHandler)

	// Register handler
	http.Handle(subscriptionsHandler.PathPrefix(), subscriptionsHandlerWithAuth)

	// Start HTTP server
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8083"
	}
	log.Printf("Starting subscriptions service HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
