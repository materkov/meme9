package controller

import (
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web2/lib"
	"github.com/materkov/meme9/web2/store"
	"github.com/materkov/meme9/web2/types"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Store *store.Store
}

func (s *Server) Serve() {
	http.HandleFunc("/feed", s.routingWrapper(s.handleFeed))
	http.HandleFunc("/posts/", s.routingWrapper(s.handlePostPage))
	http.HandleFunc("/users/", s.routingWrapper(s.handleUserPage))
	http.HandleFunc("/add_post", s.routingWrapper(s.handleAddPost))
	http.HandleFunc("/new_post", s.routingWrapper(s.handleNewPost))
	http.HandleFunc("/vk", s.routingWrapper(s.handleVkAuth))
	http.HandleFunc("/vk-callback", s.handleVkAuthCallback)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../front2/dist"))))

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

	// preload users
	userIds := lib.IdsSet{}
	for _, post := range posts {
		userIds.Add(post.UserID)
	}

	users, err := s.Store.Object.GetUsers(userIds.Get())
	//users, err := s.Store.User.GetByIdMany(userIds.Get())
	if err != nil {
		log.Printf("Error getting users: %s", err)
	}

	renderer := &types.FeedRenderer{
		Posts: make([]*types.PostRenderer, len(posts)),
	}

	for i, post := range posts {
		user := users[post.UserID]

		userName := fmt.Sprintf("User #%d", post.UserID)
		if user != nil {
			userName = user.Name
		}

		renderer.Posts[i] = &types.PostRenderer{
			ID:         strconv.Itoa(post.ID),
			AuthorName: userName,
			AuthorHref: fmt.Sprintf("/users/%d", post.UserID),
			Text:       post.Text,
		}
	}

	data := types.UniversalRenderer{
		FeedRenderer: renderer,
	}
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) handlePostPage(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	defer func() {
		log.Printf("Time: %d ms", time.Since(started).Milliseconds())
	}()

	postIDStr := strings.TrimLeft(r.URL.Path, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	//post, err := s.Store.Post.GetById(postID)
	post, err := s.Store.Object.GetPost(postID)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	} else if post == nil {
		fmt.Fprintf(w, "post not found")
		return
	}

	//user, err := s.Store.User.GetById(post.UserID)
	user, err := s.Store.Object.GetUser(post.UserID)
	if err != nil {
		log.Printf("Error getting user: %s", err)
	}

	userName := fmt.Sprintf("User #%d", post.UserID)
	if user != nil {
		userName = user.Name
	}

	renderer := &types.PostRenderer{
		ID:         strconv.Itoa(post.ID),
		AuthorName: userName,
		AuthorHref: fmt.Sprintf("/users/%d", post.UserID),
		Text:       post.Text,
	}

	_ = json.NewEncoder(w).Encode(types.UniversalRenderer{PostRenderer: renderer})
}

func (s *Server) handleUserPage(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	defer func() {
		log.Printf("Time: %d ms", time.Since(started).Milliseconds())
	}()

	userIDStr := strings.TrimLeft(r.URL.Path, "/users/")
	userID, _ := strconv.Atoi(userIDStr)

	//user, err := s.Store.User.GetById(userID)
	user, err := s.Store.Object.GetUser(userID)
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

	postRenderers := make([]*types.PostRenderer, len(posts))
	for i, post := range posts {
		postRenderers[i] = &types.PostRenderer{
			ID:         strconv.Itoa(post.ID),
			AuthorName: user.Name,
			AuthorHref: fmt.Sprintf("/users/%d", user.ID),
			Text:       post.Text,
		}
	}

	renderer := &types.UserPageRenderer{
		UserName: user.Name,
		UserID:   strconv.Itoa(user.ID),
		Posts:    postRenderers,
	}

	_ = json.NewEncoder(w).Encode(types.UniversalRenderer{
		UserPageRenderer: renderer,
	})
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

	tokenCookie, _ := r.Cookie("access_token")
	if tokenCookie == nil || tokenCookie.Value == "" {
		fmt.Fprintf(w, "no auth")
		return
	}

	userID, err := lib.ParseAuthToken(tokenCookie.Value)
	if err != nil {
		fmt.Fprintf(w, "no auth")
		return
	}

	req := addPostReq{}
	err = json.NewDecoder(r.Body).Decode(&req)
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
		UserID: userID,
	}

	//err = s.Store.Post.Add(&post)
	err = s.Store.Object.Add(store.ObjectTypePost, &post)
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
	renderer := &types.NewPostRenderer{SendLabel: "Отправить"}
	_ = json.NewEncoder(w).Encode(types.UniversalRenderer{NewPostRenderer: renderer})
}

func (s *Server) handleVkAuth(w http.ResponseWriter, r *http.Request) {
	requestScheme := lib.DefaultConfig.RequestScheme
	requestHost := lib.DefaultConfig.RequestHost
	vkAppID := lib.DefaultConfig.VkAppID
	redirectURL := fmt.Sprintf("%s://%s/vk-callback", requestScheme, requestHost)
	redirectURL = url.QueryEscape(redirectURL)
	vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", vkAppID, redirectURL)

	renderer := &types.VkAuthRenderer{URL: vkURL}
	_ = json.NewEncoder(w).Encode(types.UniversalRenderer{VkAuthRenderer: renderer})
}

const pageBegin = `
<!DOCTYPE html>
<html>
<head>
    <title>My React App</title>
</head>
<body>
<div id="root"></div>
<script>
window.__initialData = 
`

const pageEnd = `
;</script>
<script src="/static/index.js"></script>
</body>
</html>
`

func (s *Server) routingWrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		isAjax := r.Header.Get("x-ajax") == "1"
		if !isAjax {
			fmt.Fprintf(w, pageBegin)
		}

		next(w, r)

		if !isAjax {
			fmt.Fprintf(w, pageEnd)
		}
	}
}

func (s *Server) handleVkAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	vkID, err := lib.ProcessVKkCallback(code)
	if err != nil {
		log.Printf("Error processing VK: %s", err)
		fmt.Fprintf(w, "error")
		return
	}

	user, err := lib.GetOrCreateUserByVkID(s.Store, vkID)
	if err != nil {
		log.Printf("Error getting VK user: %s", err)
		fmt.Fprintf(w, "error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    lib.GenerateAuthToken(user.ID),
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/feed", 302)
}
