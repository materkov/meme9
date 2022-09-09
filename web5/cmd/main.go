package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var redisClient *redis.Client

type Post struct {
	ID     int
	Date   int
	Text   string
	UserID int
}

type User struct {
	ID   int
	Name string
	VkID int
}

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

func handleFeed(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value("viewer").(*Viewer)

	postIdsStr, err := redisClient.LRange(context.Background(), "feed", 0, 10).Result()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
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
	_ = json.NewEncoder(w).Encode(resp)
}

func handleAddPost(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	if text == "" {
		fmt.Fprintf(w, "Empty text")
		return
	}

	viewer := r.Context().Value("viewer").(*Viewer)

	nextId := int(time.Now().UnixMilli())

	post := Post{
		ID:     nextId,
		Text:   text,
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
	}
	postBytes, _ := json.Marshal(post)
	_, err := redisClient.Set(context.Background(), fmt.Sprintf("node:%d", post.ID), postBytes, 0).Result()
	if err != nil {
		fmt.Fprintf(w, "Error saving post")
		return
	}

	_, err = redisClient.LPush(context.Background(), "feed", post.ID).Result()
	if err != nil {
		log.Printf("Error saving post to feed")
	}

	_, err = redisClient.LPush(context.Background(), fmt.Sprintf("feed:%d", post.UserID), post.ID).Result()
	if err != nil {
		log.Printf("Error saving user feed key: %s", err)
	}

	apiPost := &ApiPost{
		ID:   strconv.Itoa(post.ID),
		Text: post.Text,
	}

	json.NewEncoder(w).Encode(apiPost)
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")

	users := usersList([]string{userID})
	if len(users) == 0 {
		fmt.Fprintf(w, "User not found")
		return
	}

	postIdsStr, err := redisClient.LRange(context.Background(), fmt.Sprintf("feed:%s", userID), 0, 10).Result()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	users[0].Posts = postsList(postIdsStr)

	_ = json.NewEncoder(w).Encode(users[0])
}

func handlePostPage(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("id")

	posts := postsList([]string{postID})
	if len(posts) == 0 {
		fmt.Fprintf(w, "Post not found")
		return
	}

	users := usersList([]string{posts[0].UserID})
	posts[0].User = users[0]

	_ = json.NewEncoder(w).Encode(posts[0])
}

func usersList(ids []string) []*ApiUser {
	keys := make([]string, len(ids))
	for i, userID := range ids {
		keys[i] = fmt.Sprintf("node:%s", userID)
	}

	userBytesList, err := redisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		log.Printf("Error getting users: %s", err)
	}

	var users []*User
	for _, userBytes := range userBytesList {
		if userBytes == nil {
			continue
		}

		user := &User{}
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

	postsBytes, err := redisClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		log.Printf("error getting posts: %s", err)
	}

	var posts []*Post
	for _, postBytes := range postsBytes {
		if postBytes == nil {
			continue
		}

		post := &Post{}
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
		fmt.Fprintf(w, "error")
		return
	}

	userID, err := usersGetOrCreateByVKID(vkID)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}

	authToken, err := authCreateToken(userID)

	resp := []interface{}{
		authToken,
		userID,
	}
	_ = json.NewEncoder(w).Encode(resp)
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
	_ = json.NewEncoder(w).Encode(resp)
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

		authToken := r.Header.Get("authorization")
		authToken = strings.TrimPrefix(authToken, "Bearer ")
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
		_ = json.Unmarshal(dat, &DefaultConfig)
	}

	config := os.Getenv("CONFIG")
	if config != "" {
		_ = json.Unmarshal([]byte(config), &DefaultConfig)
	}

	http.HandleFunc("/api/feed", wrapper(handleFeed))
	http.HandleFunc("/api/addPost", wrapper(handleAddPost))
	http.HandleFunc("/api/userPage", wrapper(handleUserPage))
	http.HandleFunc("/api/postPage", wrapper(handlePostPage))
	http.HandleFunc("/api/vkCallback", wrapper(handleVkCallback))
	http.HandleFunc("/api/viewer", wrapper(handleViewer))

	redisClient = redis.NewClient(&redis.Options{})

	http.ListenAndServe("127.0.0.1:8000", nil)
}
