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
	"github.com/materkov/meme9/web5/pkg/metrics"
	"github.com/materkov/meme9/web5/pkg/posts"
	"github.com/materkov/meme9/web5/pkg/telegram"
	"github.com/materkov/meme9/web5/pkg/utils"
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
			_, _ = fmt.Fprint(w, err.Error())
		} else if err != nil {
			log.Printf("[ERROR] Internal error: %s", err)
			_, _ = fmt.Fprint(w, "internal error")
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
	viewer := r.Context().Value(ViewerKey).(*Viewer)

	offset, _ := strconv.Atoi(r.FormValue("cursor"))
	limit := 10

	pipe := store.RedisClient.Pipeline()

	lenCmd := pipe.LLen(r.Context(), "feed")
	rangeCmd := pipe.LRange(r.Context(), "feed", int64(offset), int64(offset+limit-1))

	_, err := pipe.Exec(r.Context())
	if err != nil {
		write(w, nil, err)
		return
	}

	feedLen := int(lenCmd.Val())

	nextCursor := ""
	if offset+limit < feedLen {
		nextCursor = strconv.Itoa(offset + limit)
	}

	apiPosts := postsList(r.Context(), parseIds(rangeCmd.Val()), viewer.UserID)

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

	resp := struct {
		ViewerID   string  `json:"viewerId"`
		Posts      []*Post `json:"posts"`
		NextCursor string  `json:"nextCursor"`
	}{
		ViewerID:   viewerID,
		Posts:      apiPosts,
		NextCursor: nextCursor,
	}
	write(w, resp, nil)
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

	posts := postsList(r.Context(), []int{postID}, viewer.UserID)
	write(w, posts[0], nil)
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))

	viewer := r.Context().Value(ViewerKey).(*Viewer)

	users := usersList([]int{userID}, viewer.UserID, true, true)
	if len(users) == 0 {
		write(w, nil, ApiError("user not found"))
		return
	}
	users[0].Posts = userPagePosts(r.Context(), userID, 0, viewer.UserID)

	write(w, []interface{}{
		users[0],
		viewer.GetUserIDStr(),
	}, nil)
}

func handleUserPopup(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))

	userChan := make(chan store.User)
	go func() {
		user := store.User{}
		err := store.NodeGet(userID, &user)
		if err != nil {
			log.Printf("[ERROR] Error getting user: %s", err)
		}

		userChan <- user
	}()

	postsCount := make(chan int)
	go func() {
		count, err := posts.FeedLen(userID)
		utils.LogIfErr(err)
		postsCount <- count
	}()

	write(w, []interface{}{
		(<-userChan).Name,
		<-postsCount,
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
	viewer := r.Context().Value(ViewerKey).(*Viewer)

	result := userPagePosts(r.Context(), userID, cursor.Offset, viewer.UserID)

	write(w, []interface{}{
		result,
	}, nil)
}

func userPagePosts(ctx context.Context, userID int, offset int, viewerID int) *UserPostsConnection {
	redisKey := fmt.Sprintf("feed:%d", userID)
	pipe := store.RedisClient.Pipeline()

	postsIdsCmd := pipe.LRange(ctx, redisKey, int64(offset), int64(offset+10-1))
	lenCmd := pipe.LLen(ctx, redisKey)
	_, err := pipe.Exec(ctx)
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
		Items:      postsList(ctx, parseIds(postIdsStr), viewerID),
		NextCursor: nextCursor,
	}
}

func handleUserEdit(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))

	var user *store.User
	err := store.NodeGet(userID, user)
	if err != nil {
		write(w, nil, err)
		return
	} else if user == nil {
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

func handlePostPage(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("id")
	viewer := r.Context().Value(ViewerKey).(*Viewer)

	posts := postsList(r.Context(), parseIds([]string{postID}), viewer.UserID)
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

func handleViewer(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value(ViewerKey).(*Viewer)

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

func handleExecute(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value(ViewerKey).(*Viewer)

	type operation struct {
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
	}
	var operations []operation
	_ = json.Unmarshal([]byte(r.FormValue("operations")), &operations)

	var allResult []interface{}
	for _, op := range operations {
		var result interface{}

		switch op.Method {
		case "userPostsCount":
			type params struct {
				UserID string `json:"userId"`
			}
			type response struct {
				PostsCount int `json:"postsCount,omitempty"`
			}
			p := params{}
			_ = json.Unmarshal(op.Params, &p)

			userID, _ := strconv.Atoi(p.UserID)
			postsCount, _ := posts.FeedLen(userID)

			result = response{PostsCount: postsCount}

		case "userDetails":
			type params struct {
				UserIds []string `json:"userIds"`
			}
			type response struct {
				Users []*User `json:"users,omitempty"`
			}

			p := params{}
			_ = json.Unmarshal(op.Params, &p)

			users := usersList(parseIds(p.UserIds), viewer.UserID, false, false)

			result = response{Users: users}
		}

		allResult = append(allResult, result)
	}

	_ = json.NewEncoder(w).Encode(allResult)
}

type Viewer struct {
	UserID int

	Origin    string
	RequestID int
}

func (v *Viewer) GetUserIDStr() string {
	if v.UserID == 0 {
		return ""
	}
	return strconv.Itoa(v.UserID)
}

type contextKey int

const (
	ViewerKey contextKey = iota
	RequestID
	StartTime
)

func wrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}

		requestID := int(rand.Int63())

		start := time.Now()
		defer func() {
			t := time.Since(start)
			go func() {
				_ = metrics.WriteSpan(r.URL.Path, requestID, true, t)
			}()
		}()

		authToken := r.FormValue("token")
		userID, _ := authCheckToken(authToken)

		viewer := &Viewer{
			UserID:    userID,
			Origin:    r.Header.Get("origin"),
			RequestID: requestID,
		}

		ctx := context.WithValue(r.Context(), ViewerKey, viewer)
		ctx = context.WithValue(ctx, RequestID, viewer.RequestID)

		r = r.WithContext(ctx)

		next(w, r)
	}
}

type redisHook struct {
}

func (h redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, StartTime, time.Now()), nil
}

func (h redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	startTime := ctx.Value(StartTime).(time.Time)
	requestID, ok := ctx.Value(RequestID).(int)
	if ok {
		_ = metrics.WriteSpan(fmt.Sprintf("REDIS %s", cmd.Name()), requestID, false, time.Since(startTime))
	}
	return nil
}

func (h redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, StartTime, time.Now()), nil
}

func (h redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	startTime := ctx.Value(StartTime).(time.Time)
	requestID, ok := ctx.Value(RequestID).(int)
	if ok {
		for _, cmd := range cmds {
			_ = metrics.WriteSpan(fmt.Sprintf("REDIS %s", cmd.Name()), requestID, false, time.Since(startTime))
		}
	}
	return nil
}

func main() {
	queue := ""
	flag.StringVar(&queue, "queue", "", "Queue listen to")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	store.RedisClient = redis.NewClient(&redis.Options{})

	store.RedisClient.AddHook(&redisHook{})

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
	http.HandleFunc("/api/userPopup", wrapper(handleUserPopup))
	http.HandleFunc("/api/userEdit", wrapper(handleUserEdit))
	http.HandleFunc("/api/userFollow", wrapper(handleUserFollow))
	http.HandleFunc("/api/userUnfollow", wrapper(handleUserUnfollow))
	http.HandleFunc("/api/postPage", wrapper(handlePostPage))
	http.HandleFunc("/api/postDelete", wrapper(handlePostDelete))
	http.HandleFunc("/api/postLike", wrapper(handlePostLike))
	http.HandleFunc("/api/postUnlike", wrapper(handlePostUnlike))
	http.HandleFunc("/api/vkCallback", wrapper(handleVkCallback))
	http.HandleFunc("/api/viewer", wrapper(handleViewer))
	http.HandleFunc("/api/emailRegister", wrapper(handleEmailRegister))
	http.HandleFunc("/api/emailLogin", wrapper(handleAuthEmail))
	http.HandleFunc("/api/uploadAvatar", wrapper(handleUploadAvatar))
	http.HandleFunc("/api/execute", wrapper(handleExecute))

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
