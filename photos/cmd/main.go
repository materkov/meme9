package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/materkov/meme9/photos/api"
	"github.com/materkov/meme9/photos/auth"
	authpb "github.com/materkov/meme9/photos/internal/authclient/pb/github.com/materkov/meme9/api/auth"
	"github.com/materkov/meme9/photos/processor"
	"github.com/materkov/meme9/photos/uploader"
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

	authClient := authpb.NewAuthProtobufClient(
		os.Getenv("AUTH_SERVICE"),
		&http.Client{
			Timeout: 10 * time.Second,
		},
	)
	authService := auth.New(authClient)

	srv := api.New(proc, up, authService)
	_ = srv.Start()
}
