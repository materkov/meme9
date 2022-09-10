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

type Post struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Date   string `json:"date"`
	UserID string `json:"userId"`
	User   *User  `json:"user"`
}

type User struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Posts []*Post `json:"posts"`
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

func parseIds(idsStr []string) []int {
	result := make([]int, len(idsStr))
	for i, idStr := range idsStr {
		result[i], _ = strconv.Atoi(idStr)
	}

	return result
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value("viewer").(*Viewer)

	postIdsStr, err := store.RedisClient.LRange(context.Background(), "feed", 0, 10).Result()
	if err != nil {
		write(w, nil, err)
		return
	}

	apiPosts := postsList(parseIds(postIdsStr))

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
	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	}

	postID, err := postsAdd(text, viewer.UserID)
	if err != nil {
		write(w, nil, err)
		return
	}

	posts := postsList([]int{postID})
	write(w, posts[0], nil)
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")

	viewer := r.Context().Value("viewer").(*Viewer)

	users := usersList([]string{userID})
	if len(users) == 0 {
		write(w, nil, ApiError("user not found"))
		return
	}

	postIdsStr, err := store.RedisClient.LRange(context.Background(), fmt.Sprintf("feed:%s", userID), 0, 10).Result()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	users[0].Posts = postsList(parseIds(postIdsStr))

	write(w, []interface{}{
		users[0],
		viewer.GetUserIDStr(),
	}, nil)
}

func handleUserEdit(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))

	user := &store.User{}
	err := store.NodeGet(userID, user)
	if err != nil {
		write(w, nil, err)
		return
	} else if user == nil {
		write(w, nil, ApiError("user not found"))
		return
	}

	viewer := r.Context().Value("viewer").(*Viewer)
	if viewer.UserID != user.ID {
		write(w, nil, ApiError("no access to edit this user"))
		return
	}

	name := r.FormValue("name")
	if name == "" {
		write(w, nil, ApiError("name is empty"))
		return
	} else if len(name) > 100 {
		write(w, nil, ApiError("name is too long"))
		return
	}

	user.Name = name

	err = store.NodeSave(user.ID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	write(w, nil, nil)
}

func handlePostPage(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("id")

	posts := postsList(parseIds([]string{postID}))
	if len(posts) == 0 {
		write(w, nil, ApiError("post not found"))
		return
	}

	users := usersList([]string{posts[0].UserID})
	posts[0].User = users[0]

	write(w, posts[0], nil)
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
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := []interface{}{
		authToken,
		userID,
	}
	write(w, resp, nil)
}

func handleViewer(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value("viewer").(*Viewer)

	var user *User
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

func (v *Viewer) GetUserIDStr() string {
	if v.UserID == 0 {
		return ""
	}
	return strconv.Itoa(v.UserID)
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
	http.HandleFunc("/api/userEdit", wrapper(handleUserEdit))
	http.HandleFunc("/api/postPage", wrapper(handlePostPage))
	http.HandleFunc("/api/vkCallback", wrapper(handleVkCallback))
	http.HandleFunc("/api/viewer", wrapper(handleViewer))

	store.RedisClient = redis.NewClient(&redis.Options{})

	http.ListenAndServe("127.0.0.1:8000", nil)
}
