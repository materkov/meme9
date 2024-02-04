package main

import "github.com/materkov/meme9/rss2/pb/github.com/materkov/meme9/api"

func main() {
	trp := api.NewPostsProtobufClient("http://localhost:8000", nil)
	trp.List(nil, nil)
}
