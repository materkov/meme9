package controller

import (
	"fmt"
	"github.com/materkov/meme9/web2/lib"
	"github.com/materkov/meme9/web2/store"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Store *store.Store
}

func (s *Server) Serve() {
	http.HandleFunc("/feed", s.handleFeed)

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	defer func() {
		log.Printf("Time: %d ms", time.Since(started).Milliseconds())
	}()

	posts, err := s.Store.Post.GetAll()
	if err != nil {
		log.Printf("%s", err)
	}

	// Preload users
	userIds := lib.IdsSet{}
	for _, post := range posts {
		userIds.Add(post.UserID)
	}

	users, err := s.Store.User.GetByIdMany(userIds.Get())
	if err != nil {
		log.Printf("Error getting users: %s", err)
	}

	postsHTML := ""
	for _, post := range posts {
		user := users[post.UserID]

		renderer := PostRenderer{
			post: &post,
			user: user,
		}

		postsHTML += renderer.Render()
	}

	fmt.Fprintf(w, "<html><body><h1>Feed:</h1>")
	fmt.Fprintf(w, "%s", postsHTML)
	fmt.Fprintf(w, "</body></html>")
}
