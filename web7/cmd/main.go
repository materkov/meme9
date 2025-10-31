package main

import (
	"context"
	"log"
	"os"

	"github.com/materkov/meme9/web7/adapters/mongo"
	"github.com/materkov/meme9/web7/api"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@localhost:27017/meme9?authSource=admin"
	}

	ctx := context.Background()
	mongoAdapter, err := mongo.NewAdapter(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping MongoDB to verify connection
	err = mongoAdapter.Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")

	apiAdapter := api.NewAPI(mongoAdapter)
	apiAdapter.Serve()
}
