package main

import (
	"context"
	"log"
	"net/http"
	"os"

	usersapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/users"
	"github.com/materkov/meme9/users-service/adapters/users"
	"github.com/materkov/meme9/users-service/api"
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

	// Ensure indexes
	err = usersAdapter.EnsureIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to ensure user indexes: %v", err)
	}

	// Initialize services
	usersService := api.NewService(usersAdapter)

	// Create Twirp server
	usersHandler := usersapi.NewUsersServer(usersService)
	usersHandlerWithCORS := api.AuthMiddleware(api.CORSMiddleware(usersHandler))

	// Register handler
	http.Handle(usersHandler.PathPrefix(), usersHandlerWithCORS)

	// Start HTTP server
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8082"
	}
	log.Printf("Starting users service HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
