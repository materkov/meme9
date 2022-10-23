package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/api"
	"github.com/materkov/meme9/web5/pkg/telegram"
	"github.com/materkov/meme9/web5/store"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Date   string `json:"date"`
	UserID string `json:"userId"`
	User   *User  `json:"user"`

	CanDelete bool `json:"canDelete,omitempty"`

	LikesCount int  `json:"likesCount,omitempty"`
	IsLiked    bool `json:"isLiked,omitempty"`
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`

	IsFollowing bool `json:"isFollowing,omitempty"`
}

type ApiError string

func (e ApiError) Error() string {
	return string(e)
}

func write(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)

		var apiErr ApiError
		if errors.As(err, &apiErr) {
			_, _ = fmt.Fprint(w, err.Error())
		} else if err != nil {
			log.Printf("[ERROR] Internal error: %s", err)
			_, _ = fmt.Fprint(w, "internal error")
		}
	} else {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func handleAddPost(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	if text == "" {
		write(w, nil, ApiError("text is empty"))
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	}

	postID, err := postsAdd(text, viewer.UserID)
	if err != nil {
		write(w, nil, err)
		return
	}

	posts := postsList([]int{postID}, viewer.UserID)
	write(w, posts[0], nil)
}

func handleUserEdit(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))

	user := store.User{}
	err := store.NodeGet(userID, &user)
	if err != nil {
		write(w, nil, err)
		return
	} else if user.ID == 0 {
		write(w, nil, ApiError("user not found"))
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
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

func handleUserFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))
	viewer := r.Context().Value(ViewerKey).(*Viewer)

	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	} else if userID == 0 {
		write(w, nil, ApiError("empty user"))
		return
	} else if userID == viewer.UserID {
		write(w, nil, ApiError("you cannot subscribe to yourself"))
		return
	}

	err := usersFollow(viewer.UserID, userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	write(w, nil, nil)
}

func handleUserUnfollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))
	viewer := r.Context().Value(ViewerKey).(*Viewer)

	err := usersUnfollow(viewer.UserID, userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	write(w, nil, nil)
}

func handlePostDelete(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.FormValue("id"))

	post := &store.Post{}
	err := store.NodeGet(postID, post)
	if err == store.ErrNodeNotFound {
		write(w, nil, ApiError("post not found"))
		return
	} else if err != nil {
		write(w, nil, err)
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if post.UserID != viewer.UserID {
		write(w, nil, ApiError("no access to delete this post"))
		return
	}

	if !post.IsDeleted {
		err = postsDelete(post)
		if err != nil {
			write(w, nil, err)
			return
		}
	}

	write(w, []interface{}{}, nil)
}

func handlePostLike(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.FormValue("id"))

	post := store.Post{}
	err := store.NodeGet(postID, &post)
	if err == store.ErrNodeNotFound {
		write(w, nil, ApiError("post not found"))
		return
	} else if err != nil {
		write(w, nil, err)
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	}

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("postLikes:%d", postID)
	_ = pipe.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: viewer.UserID,
	})
	cardCmd := pipe.ZCard(context.Background(), key)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := struct {
		LikesCount int `json:"likesCount"`
	}{
		LikesCount: int(cardCmd.Val()),
	}

	write(w, resp, nil)
}

func handlePostUnlike(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.FormValue("id"))

	post := store.Post{}
	err := store.NodeGet(postID, &post)
	if err == store.ErrNodeNotFound {
		write(w, nil, ApiError("post not found"))
		return
	} else if err != nil {
		write(w, nil, err)
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	}

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("postLikes:%d", postID)
	pipe.ZRem(context.Background(), key, viewer.UserID)
	cardCmd := pipe.ZCard(context.Background(), key)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := struct {
		LikesCount int `json:"likesCount"`
	}{
		LikesCount: int(cardCmd.Val()),
	}

	write(w, resp, nil)
}

func handleVkCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirectUri")

	vkID, vkAccessToken, err := authExchangeCode(code, redirectURI)
	if err != nil {
		write(w, nil, err)
		return
	}

	userID, err := usersGetOrCreateByVKID(vkID)
	if err != nil {
		write(w, nil, err)
		return
	}

	user := &store.User{}
	err = store.NodeGet(userID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	user.VkAccessToken = vkAccessToken
	err = store.NodeSave(user.ID, user)
	if err != nil {
		log.Printf("error saving user")
	}

	_, _ = store.RedisClient.RPush(context.Background(), "queue", user.ID).Result()

	authToken, err := authCreateToken(userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	err = telegram.SendNotify(fmt.Sprintf("meme new login: https://vk.com/id%d", user.VkID))
	if err != nil {
		log.Printf("Error sending telegram notify: %s", err)
	}

	resp := []interface{}{
		authToken,
		userID,
	}
	write(w, resp, nil)
}

func handleEmailRegister(w http.ResponseWriter, r *http.Request) {
	email, password := r.FormValue("email"), r.FormValue("password")
	validateErr := authValidateCredentials(email, password)
	if validateErr != "" {
		write(w, nil, ApiError(validateErr))
		return
	}

	userID, err := authRegister(email, password)
	if err != nil {
		write(w, nil, err)
		return
	}

	users := usersList([]int{userID}, userID, false, false)

	token, err := authCreateToken(userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := struct {
		Token string `json:"token"`
		User  *User  `json:"user"`
	}{
		Token: token,
		User:  users[0],
	}
	write(w, resp, nil)
}

func handleAuthEmail(w http.ResponseWriter, r *http.Request) {
	userID, err := authEmailAuth(r.FormValue("email"), r.FormValue("password"))
	if err == ErrInvalidCredentials {
		write(w, nil, ApiError("invalid credentials"))
		return
	} else if err != nil {
		write(w, nil, err)
		return
	}

	token, err := authCreateToken(userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	users := usersList([]int{userID}, userID, false, false)

	resp := struct {
		User  *User  `json:"user"`
		Token string `json:"token"`
	}{
		User:  users[0],
		Token: token,
	}
	write(w, resp, nil)
}

func handleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		write(w, nil, ApiError("invalid file"))
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		write(w, nil, err)
		return
	}

	hash := sha256.Sum256(fileBytes)
	hashHex := hex.EncodeToString(hash[:])

	err = filesSelectelUpload(fileBytes, hashHex)
	if err != nil {
		write(w, nil, err)
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	}

	user := &store.User{}
	err = store.NodeGet(viewer.UserID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	user.AvatarSha = hashHex

	err = store.NodeSave(user.ID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := struct {
		Avatar string `json:"avatar"`
	}{
		Avatar: filesGetURL(hashHex),
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

type contextKey int

const ViewerKey contextKey = iota

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
		r = r.WithContext(context.WithValue(r.Context(), ViewerKey, viewer))

		next(w, r)
	}
}

func main() {
	queue := ""
	flag.StringVar(&queue, "queue", "", "Queue listen to")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	store.RedisClient = redis.NewClient(&redis.Options{})

	configStr, err := store.RedisClient.Get(context.Background(), "config").Bytes()
	if err != nil {
		log.Fatalf("Failed reading config: %s", err)
	}

	err = json.Unmarshal(configStr, &store.DefaultConfig)
	if err != nil {
		log.Fatalf("Error parsing config JSON: %s", err)
	}

	if queue != "" {
		HandleWorker(queue)
		return
	}

	http.HandleFunc("/api", api.HandleAPI)

	http.HandleFunc("/api/addPost", wrapper(handleAddPost))
	http.HandleFunc("/api/userEdit", wrapper(handleUserEdit))
	http.HandleFunc("/api/userFollow", wrapper(handleUserFollow))
	http.HandleFunc("/api/userUnfollow", wrapper(handleUserUnfollow))
	http.HandleFunc("/api/postDelete", wrapper(handlePostDelete))
	http.HandleFunc("/api/postLike", wrapper(handlePostLike))
	http.HandleFunc("/api/postUnlike", wrapper(handlePostUnlike))
	http.HandleFunc("/api/vkCallback", wrapper(handleVkCallback))
	http.HandleFunc("/api/emailRegister", wrapper(handleEmailRegister))
	http.HandleFunc("/api/emailLogin", wrapper(handleAuthEmail))
	http.HandleFunc("/api/uploadAvatar", wrapper(handleUploadAvatar))

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func HandleWorker(queue string) {
	for {
		result, err := store.RedisClient.BLPop(context.Background(), time.Second*5, queue).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return
		}

		log.Printf("Got queue task: %v", result)

		userID, _ := strconv.Atoi(result[1])
		err = usersRefreshFromVk(userID)
		if err != nil {
			log.Printf("Error doing queue: %s", err)
		}
	}
}
