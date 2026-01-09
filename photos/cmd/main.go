package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/materkov/meme9/photos/api"
	"github.com/materkov/meme9/photos/auth"
	authpb "github.com/materkov/meme9/photos/internal/authclient/pb/github.com/materkov/meme9/api/auth"
	"github.com/materkov/meme9/photos/processor"
	"github.com/materkov/meme9/photos/uploader"
	photosapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/photos"
)

func main() {
	_ = godotenv.Load()

	proc := processor.New()

	up, err := uploader.New(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_KEY"),
	)
	if err != nil {
		panic(err)
	}

	authServiceURL := os.Getenv("AUTH_SERVICE")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8081"
	}

	authClient := authpb.NewAuthProtobufClient(
		authServiceURL,
		&http.Client{
			Timeout: 10 * time.Second,
		},
	)
	authService := auth.New(authClient)

	// Create API service for upload endpoint
	apiService := api.New(proc, up, authService)

	// Create Twirp service
	photosService := api.NewPhotos()
	photosHandler := photosapi.NewPhotosServer(photosService)
	photosHandlerWithCORS := api.CORSMiddleware(photosHandler)

	// Setup routes
	mux := http.NewServeMux()
	// Upload endpoint (handles file uploads)
	mux.Handle("/twirp/meme.photos.Photos/upload", apiService.Routes())
	// Twirp service endpoints
	mux.Handle(photosHandler.PathPrefix(), photosHandlerWithCORS)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8086"
	}

	log.Printf("Starting photos service HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
