package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/materkov/meme9/web2/store"
)

type Server struct {
	Store *store.Store
}

func (s *Server) Serve() {
	http.HandleFunc("/feed", s.handleFeed)

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	posts, err := s.Store.Post.GetAll()
	if err != nil {
		log.Printf("%s", err)
	}

	postsHTML := ""
	for _, post := range posts {
		renderer := PostRenderer{
			post: &post,
			user: &store.User{ID: post.UserID},
		}

		postsHTML += renderer.Render()
	}

	fmt.Fprintf(w, "<html><body><h1>Feed:</h1>")
	fmt.Fprintf(w, "%s", postsHTML)
	fmt.Fprintf(w, "</body></html>")
}
