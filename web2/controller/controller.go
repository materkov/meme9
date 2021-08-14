package controller

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/materkov/meme9/web2/lib"
	"github.com/materkov/meme9/web2/store"
	"io/ioutil"
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
	http.HandleFunc("/feed", s.handleFeed)
	http.HandleFunc("/posts/", s.handlePostPage)
	http.HandleFunc("/users/", s.handleUserPage)
	http.HandleFunc("/add_post", s.handleAddPost)
	http.HandleFunc("/new_post", s.handleNewPost)
	http.HandleFunc("/vk", s.handleVkAuth)
	http.HandleFunc("/vk-callback", s.handleVkAuthCallback)

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

type MyCustomClaims struct {
	jwt.StandardClaims
	UserID int `json:"userId"`
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

	token, err := jwt.ParseWithClaims(tokenCookie.Value, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(lib.DefaultConfig.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		fmt.Fprintf(w, "no auth")
		return
	}
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		fmt.Fprintf(w, "no auth")
		return
	}

	userID := token.Claims.(*MyCustomClaims).UserID

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

func (s *Server) handleVkAuth(w http.ResponseWriter, r *http.Request) {
	renderer := VkAuthRenderer{}
	fmt.Fprint(w, renderer.Render())
}

func (s *Server) handleVkAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	redirectURI := fmt.Sprintf("%s://%s/vk-callback", lib.DefaultConfig.RequestScheme, lib.DefaultConfig.RequestHost)

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(lib.DefaultConfig.VkAppID)},
		"client_secret": []string{lib.DefaultConfig.VkAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		fmt.Fprintf(w, "http error")
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "http error")
		return
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		fmt.Fprintf(w, "http error")
		return
	} else if body.AccessToken == "" {
		fmt.Fprintf(w, "http error")
		return
	}

	user, err := s.Store.User.GetByVkID(body.UserID)
	if err != nil {
		fmt.Fprintf(w, "http error")
		return
	}

	if user == nil {
		user = &store.User{
			Name: fmt.Sprintf("VK User #%d", body.UserID),
			VkID: body.UserID,
		}
		err := s.Store.User.Add(user)
		if err != nil {
			fmt.Fprintf(w, "http error")
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "meme9",
		},
		UserID: user.ID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(lib.DefaultConfig.JwtSecret))

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/feed", 302)
}
