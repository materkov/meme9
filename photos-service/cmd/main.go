package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	photosapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/photos"
	"github.com/materkov/meme9/photos-service/api"
	"github.com/materkov/meme9/photos-service/processor"
	"github.com/materkov/meme9/photos-service/uploader"
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

	// Create API service for upload endpoint
	apiService := api.New(proc, up)

	// Create Twirp service
	photosService := api.NewPhotos()
	photosHandler := photosapi.NewPhotosServer(photosService)

	// Setup routes
	mux := http.NewServeMux()
	// Upload endpoint (handles file uploads) - with auth middleware
	uploadHandler := api.AuthMiddleware(apiService.Routes())
	mux.Handle("/twirp/meme.photos.Photos/upload", uploadHandler)
	// Twirp service endpoints
	mux.Handle(photosHandler.PathPrefix(), photosHandler)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8086"
	}

	log.Printf("Starting photos service HTTP server at http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
