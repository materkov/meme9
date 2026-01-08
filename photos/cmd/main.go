package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/materkov/meme9/photos/api"
	"github.com/materkov/meme9/photos/processor"
	"github.com/materkov/meme9/photos/uploader"
)

func main() {
	_ = godotenv.Load()

	processor := processor.New()
	uploader, err := uploader.New(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_KEY"),
	)
	if err != nil {
		panic(err)
	}

	api := api.New(processor, uploader)
	api.Start()
}
