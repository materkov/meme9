package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	json_api "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/json_api"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	"github.com/materkov/meme9/web7/apiwrapper"
	grpcservice "github.com/materkov/meme9/web7/grpc"
	"github.com/materkov/meme9/web7/html"
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

	apiAdapter := api.NewAPI(postsAdapter, usersAdapter, tokensAdapter, subscriptionsAdapter, postsService, tokensService)

	// Create API wrapper router
	apiRouter := apiwrapper.NewRouter(apiAdapter)
	apiRouter.RegisterRoutes()

	// Create HTML router
	htmlRouter := html.NewRouter(apiAdapter)

	// Register HTML routes
	http.HandleFunc("/users/{id}", htmlRouter.UserPageHandler)
	http.HandleFunc("/posts/{id}", htmlRouter.PostPageHandler)
	http.HandleFunc("/feed", htmlRouter.FeedPageHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	http.HandleFunc("/", htmlRouter.IndexHandler)

	// Start gRPC server
	go startGRPCServer(apiAdapter)

	// Start HTTP server
	apiRouter.StartServer("127.0.0.1:8080")
}

func startGRPCServer(apiAdapter *api.API) {
	lis, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcservice.AuthInterceptor(apiAdapter)),
	)

	grpcService := grpcservice.NewServer(apiAdapter)
	json_api.RegisterJsonAPIServer(grpcServer, grpcService)

	// Enable server reflection for grpcui
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server at %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
