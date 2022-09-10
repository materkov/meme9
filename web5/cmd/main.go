package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/store"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ApiPost struct {
	ID     string   `json:"id"`
	Text   string   `json:"text"`
	Date   string   `json:"date"`
	UserID string   `json:"userId"`
	User   *ApiUser `json:"user"`
}

type ApiUser struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Posts []*ApiPost `json:"posts"`
}

type ApiError string

func (e ApiError) Error() string {
	return string(e)
}

func write(w http.ResponseWriter, data interface{}, err error) {
	resp := struct {
		Ok    bool        `json:"ok"`
		Data  interface{} `json:"data,omitempty"`
		Error string      `json:"error,omitempty"`
	}{}

	var apiErr ApiError
	if errors.As(err, &apiErr) {
		resp.Error = err.Error()
	} else if err != nil {
		log.Printf("[ERROR] Internal error: %s", err)
		resp.Error = "internal error"
	} else {
		resp.Ok = true
		resp.Data = data
	}
	_ = json.NewEncoder(w).Encode(&resp)
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value("viewer").(*Viewer)

	postIdsStr, err := store.RedisClient.LRange(context.Background(), "feed", 0, 10).Result()
	if err != nil {
		write(w, nil, err)
		return
	}

	apiPosts := postsList(postIdsStr)

	for _, post := range apiPosts {
		users := usersList([]string{post.UserID})
		if len(users) == 1 {
			post.User = users[0]
		}
	}

	viewerID := ""
	if viewer.UserID != 0 {
		viewerID = strconv.Itoa(viewer.UserID)
	}

	resp := []interface{}{
		viewerID,
		apiPosts,
	}
	write(w, resp, nil)
}

func handleAddPost(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	if text == "" {
		write(w, nil, ApiError("text is empty"))
		return
	}

	viewer := r.Context().Value("viewer").(*Viewer)

	nextId := int(time.Now().UnixMilli())

	post := store.Post{
		ID:     nextId,
		Text:   text,
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
	}
	postBytes, _ := json.Marshal(post)
	_, err := store.RedisClient.Set(context.Background(), fmt.Sprintf("node:%d", post.ID), postBytes, 0).Result()
	if err != nil {
		write(w, nil, err)
		return
	}

	_, err = store.RedisClient.LPush(context.Background(), "feed", post.ID).Result()
	if err != nil {
		log.Printf("Error saving post to feed")
	}

	_, err = store.RedisClient.LPush(context.Background(), fmt.Sprintf("feed:%d", post.UserID), post.ID).Result()
	if err != nil {
		log.Printf("Error saving user feed key: %s", err)
	}

	apiPost := &ApiPost{
		ID:   strconv.Itoa(post.ID),
		Text: post.Text,
	}

	write(w, apiPost, nil)
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")

	users := usersList([]string{userID})
	if len(users) == 0 {
		write(w, nil, ApiError("user not found"))
		return
	}

	postIdsStr, err := store.RedisClient.LRange(context.Background(), fmt.Sprintf("feed:%s", userID), 0, 10).Result()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	users[0].Posts = postsList(postIdsStr)

	write(w, users[0], nil)
}

func handlePostPage(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("id")

	posts := postsList([]string{postID})
	if len(posts) == 0 {
		write(w, nil, ApiError("post not found"))
		return
	}

	users := usersList([]string{posts[0].UserID})
	posts[0].User = users[0]

	write(w, posts[0], nil)
}

func usersList(ids []string) []*ApiUser {
	keys := make([]string, len(ids))
	for i, userID := range ids {
		keys[i] = fmt.Sprintf("node:%s", userID)
	}

	userBytesList, err := store.RedisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		log.Printf("Error getting users: %s", err)
	}

	var users []*store.User
	for _, userBytes := range userBytesList {
		if userBytes == nil {
			continue
		}

		user := &store.User{}
		err = json.Unmarshal([]byte(userBytes.(string)), user)
		if err != nil {
			log.Printf("Error unmarshalling user: %s", err)
			continue
		}

		users = append(users, user)
	}

	apiUsers := make([]*ApiUser, len(users))
	for i, user := range users {
		apiUsers[i] = &ApiUser{
			ID:   strconv.Itoa(user.ID),
			Name: user.Name,
		}
	}

	return apiUsers
}

func postsList(ids []string) []*ApiPost {
	keys := make([]string, len(ids))
	for i, postID := range ids {
		keys[i] = fmt.Sprintf("node:%s", postID)
	}

	postsBytes, err := store.RedisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		log.Printf("error getting posts: %s", err)
	}

	var posts []*store.Post
	for _, postBytes := range postsBytes {
		if postBytes == nil {
			continue
		}

		post := &store.Post{}
		err = json.Unmarshal([]byte(postBytes.(string)), post)
		if err != nil {
			continue
		}

		posts = append(posts, post)
	}

	apiPosts := make([]*ApiPost, len(posts))
	for i, post := range posts {
		apiPost := &ApiPost{
			ID:     strconv.Itoa(post.ID),
			Text:   post.Text,
			Date:   time.Unix(int64(post.Date), 0).UTC().Format(time.RFC3339),
			UserID: strconv.Itoa(post.UserID),
		}
		apiPosts[i] = apiPost
	}

	return apiPosts
}

func handleVkCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirectUri")

	vkID, err := authExchangeCode(code, redirectURI)
	if err != nil {
		write(w, nil, err)
		return
	}

	userID, err := usersGetOrCreateByVKID(vkID)
	if err != nil {
		write(w, nil, err)
		return
	}

	authToken, err := authCreateToken(userID)

	resp := []interface{}{
		authToken,
		userID,
	}
	write(w, resp, nil)
}

func handleViewer(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value("viewer").(*Viewer)

	var user *ApiUser
	if viewer.UserID != 0 {
		users := usersList([]string{strconv.Itoa(viewer.UserID)})
		user = users[0]
	}

	resp := []interface{}{
		user,
	}
	write(w, resp, nil)
}

type Viewer struct {
	UserID int

	Origin string
}

func wrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}

		authToken := r.FormValue("token")
		userID, _ := authCheckToken(authToken)

		viewer := &Viewer{
			UserID: userID,
			Origin: r.Header.Get("origin"),
		}
		r = r.WithContext(context.WithValue(r.Context(), "viewer", viewer))

		next(w, r)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	homeDir, _ := os.UserHomeDir()
	dat, _ := os.ReadFile(homeDir + "/mypage/config.json")
	if len(dat) > 0 {
		_ = json.Unmarshal(dat, &store.DefaultConfig)
	}

	config := os.Getenv("CONFIG")
	if config != "" {
		_ = json.Unmarshal([]byte(config), &store.DefaultConfig)
	}

	http.HandleFunc("/api/feed", wrapper(handleFeed))
	http.HandleFunc("/api/addPost", wrapper(handleAddPost))
	http.HandleFunc("/api/userPage", wrapper(handleUserPage))
	http.HandleFunc("/api/postPage", wrapper(handlePostPage))
	http.HandleFunc("/api/vkCallback", wrapper(handleVkCallback))
	http.HandleFunc("/api/viewer", wrapper(handleViewer))

	store.RedisClient = redis.NewClient(&redis.Options{})

	http.ListenAndServe("127.0.0.1:8000", nil)
}
