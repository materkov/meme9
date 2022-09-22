package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/pkg/telegram"
	"github.com/materkov/meme9/web5/store"
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
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`

	Posts *UserPostsConnection `json:"posts"`

	IsFollowing bool `json:"isFollowing,omitempty"`

	FollowingCount  int `json:"followingCount,omitempty"`
	FollowedByCount int `json:"followedByCount,omitempty"`
}

type UserPostsConnection struct {
	Count      int     `json:"count,omitempty"`
	Items      []*Post `json:"items,omitempty"`
	NextCursor string  `json:"nextCursor,omitempty"`
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
			fmt.Fprintf(w, err.Error())
		} else if err != nil {
			log.Printf("[ERROR] Internal error: %s", err)
			fmt.Fprintf(w, "internal error")
		}
	} else {
		_ = json.NewEncoder(w).Encode(data)
	}
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

	apiPosts := postsList(parseIds(postIdsStr), viewer.UserID)

	for _, post := range apiPosts {
		userID, _ := strconv.Atoi(post.UserID)
		users := usersList([]int{userID}, viewer.UserID, false, false)
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

	posts := postsList([]int{postID}, viewer.UserID)
	write(w, posts[0], nil)
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))

	viewer := r.Context().Value("viewer").(*Viewer)

	users := usersList([]int{userID}, viewer.UserID, true, true)
	if len(users) == 0 {
		write(w, nil, ApiError("user not found"))
		return
	}
	users[0].Posts = userPagePosts(userID, 0, viewer.UserID)

	write(w, []interface{}{
		users[0],
		viewer.GetUserIDStr(),
	}, nil)
}

type postsListCursor struct {
	Offset int
}

func (p *postsListCursor) ToString() string {
	return strconv.Itoa(p.Offset)
}

func ParsePostsListCursor(cursor string) *postsListCursor {
	result := &postsListCursor{}
	result.Offset, _ = strconv.Atoi(cursor)
	return result
}

func handleUserPagePosts(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))
	cursor := ParsePostsListCursor(r.FormValue("cursor"))
	viewer := r.Context().Value("viewer").(*Viewer)

	result := userPagePosts(userID, cursor.Offset, viewer.UserID)

	write(w, []interface{}{
		result,
	}, nil)
}

func userPagePosts(userID int, offset int, viewerID int) *UserPostsConnection {
	redisKey := fmt.Sprintf("feed:%d", userID)
	pipe := store.RedisClient.Pipeline()

	postsIdsCmd := pipe.LRange(context.Background(), redisKey, int64(offset), int64(offset+10-1))
	lenCmd := pipe.LLen(context.Background(), redisKey)
	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	postIdsStr := postsIdsCmd.Val()
	count := int(lenCmd.Val())

	nextCursor := ""
	if offset+10 < count {
		cursor := postsListCursor{Offset: offset + 10}
		nextCursor = cursor.ToString()
	}

	return &UserPostsConnection{
		Count:      count,
		Items:      postsList(parseIds(postIdsStr), viewerID),
		NextCursor: nextCursor,
	}
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

func handleUserFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))
	viewer := r.Context().Value("viewer").(*Viewer)

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
	viewer := r.Context().Value("viewer").(*Viewer)

	err := usersUnfollow(viewer.UserID, userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	write(w, nil, nil)
}

func handlePostPage(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("id")
	viewer := r.Context().Value("viewer").(*Viewer)

	posts := postsList(parseIds([]string{postID}), viewer.UserID)
	if len(posts) == 0 {
		write(w, nil, ApiError("post not found"))
		return
	}

	userID, _ := strconv.Atoi(posts[0].UserID)
	users := usersList([]int{userID}, viewer.UserID, false, true)
	posts[0].User = users[0]

	write(w, posts[0], nil)
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

	viewer := r.Context().Value("viewer").(*Viewer)
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

func handleViewer(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value("viewer").(*Viewer)

	var user *User
	if viewer.UserID != 0 {
		users := usersList([]int{viewer.UserID}, viewer.UserID, false, true)
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

	http.HandleFunc("/api/feed", wrapper(handleFeed))
	http.HandleFunc("/api/addPost", wrapper(handleAddPost))
	http.HandleFunc("/api/userPage", wrapper(handleUserPage))
	http.HandleFunc("/api/userPage/posts", wrapper(handleUserPagePosts))
	http.HandleFunc("/api/userEdit", wrapper(handleUserEdit))
	http.HandleFunc("/api/userFollow", wrapper(handleUserFollow))
	http.HandleFunc("/api/userUnfollow", wrapper(handleUserUnfollow))
	http.HandleFunc("/api/postPage", wrapper(handlePostPage))
	http.HandleFunc("/api/postDelete", wrapper(handlePostDelete))
	http.HandleFunc("/api/vkCallback", wrapper(handleVkCallback))
	http.HandleFunc("/api/viewer", wrapper(handleViewer))

	http.ListenAndServe("127.0.0.1:8000", nil)
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
