package controller

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web2/lib"
	"github.com/materkov/meme9/web2/store"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Store *store.Store
}

func (s *Server) Serve() {
	http.HandleFunc("/feed", s.handleFeed)
	http.HandleFunc("/posts/", s.handlePostPage)
	http.HandleFunc("/users/", s.handleUserPage)
	http.HandleFunc("/add_post", s.handleAddPost)
	http.HandleFunc("/new_post", s.handleNewPost)

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

func (s *Server) handlePostPage(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	defer func() {
		log.Printf("Time: %d ms", time.Since(started).Milliseconds())
	}()

	postIDStr := strings.TrimLeft(r.URL.Path, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	post, err := s.Store.Post.GetById(postID)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	} else if post == nil {
		fmt.Fprintf(w, "post not found")
		return
	}

	user, err := s.Store.User.GetById(post.UserID)
	if err != nil {
		log.Printf("Error getting user: %s", err)
	}

	renderer := PostRenderer{
		post: post,
		user: user,
	}

	fmt.Fprintf(w, "<html><body><h1>Post page:</h1>")
	fmt.Fprintf(w, "%s", renderer.Render())
	fmt.Fprintf(w, "</body></html>")
}

func (s *Server) handleUserPage(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	defer func() {
		log.Printf("Time: %d ms", time.Since(started).Milliseconds())
	}()

	userIDStr := strings.TrimLeft(r.URL.Path, "/users/")
	userID, _ := strconv.Atoi(userIDStr)

	user, err := s.Store.User.GetById(userID)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	} else if user == nil {
		fmt.Fprintf(w, "user not found")
		return
	}

	posts, err := s.Store.Post.GetByUser(user.ID, 50)
	if err != nil {
		log.Printf("Error getting user posts: %s", err)
	}

	renderer := UserPageRenderer{
		user:  user,
		posts: posts,
	}

	fmt.Fprintf(w, "<html><body><h1>User page:</h1>")
	fmt.Fprintf(w, "%s", renderer.Render())
	fmt.Fprintf(w, "</body></html>")

}

type addPostReq struct {
	Text string
}

type addPostResp struct {
	PostID  string `json:"postId"`
	PostURL string `json:"postUrl"`
}

func (s *Server) handleAddPost(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	defer func() {
		log.Printf("Time: %d ms", time.Since(started).Milliseconds())
	}()

	req := addPostReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Fprintf(w, "err json")
		return
	}

	if req.Text == "" {
		fmt.Fprintf(w, "empty_text")
		return
	} else if len(req.Text) > 1000 {
		fmt.Fprintf(w, "text_too_long")
		return
	}

	post := store.Post{
		Text:   req.Text,
		UserID: 1,
	}

	err = s.Store.Post.Add(&post)
	if err != nil {
		log.Printf("Error adding post: %s", err)
		fmt.Fprintf(w, "internal_error")
		return
	}

	resp := addPostResp{
		PostID:  strconv.Itoa(post.ID),
		PostURL: fmt.Sprintf("/posts/%d", post.ID),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleNewPost(w http.ResponseWriter, r *http.Request) {
	renderer := NewPostRenderer{}
	fmt.Fprintf(w, renderer.Render())
}
