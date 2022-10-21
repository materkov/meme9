package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/store"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	ID     string `json:"id,omitempty"`
	URL    string `json:"url,omitempty"`
	Date   string `json:"date,omitempty"`
	Text   string `json:"text,omitempty"`
	UserID string `json:"userId,omitempty"`
}

type Edges struct {
	URL        string `json:"url,omitempty"`
	TotalCount int    `json:"totalCount,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`

	Items []string `json:"items,omitempty"`
}

type User struct {
	ID     string `json:"id,omitempty"`
	URL    string `json:"url,omitempty"`
	Name   string `json:"name,omitempty"`
	Bio    string `json:"bio,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// /feed
func handleFeed(viewerID int, _ string) []interface{} {
	postIds, err := store.RedisClient.LRange(context.Background(), "feed", 0, 20).Result()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	feed := Edges{
		URL:        "/feed",
		TotalCount: 20,
		NextCursor: "",
		Items:      postIds,
	}

	var results []interface{}
	results = append(results, feed)

	for _, postID := range postIds {
		results = append(results, handlePostsId(viewerID, "/posts/"+postID)...)
		results = append(results, handlePostsIsLiked(viewerID, "/posts/"+postID+"/isLiked")...)
	}

	return results
}

// /users/:id
func handleUserById(_ int, url string) []interface{} {
	userID, _ := strconv.Atoi(strings.TrimPrefix(url, "/users/"))

	user := store.User{}
	_ = store.NodeGet(userID, &user)

	wrapped := User{
		ID:   strconv.Itoa(userID),
		URL:  fmt.Sprintf("/users/%d", userID),
		Name: user.Name,
	}
	if user.AvatarSha != "" {
		wrapped.Avatar = files.GetURL(user.AvatarSha)
	} else if user.VkPhoto200 != "" {
		wrapped.Avatar = user.VkPhoto200
	}

	return []interface{}{wrapped}
}

// /users/:id/followers
func handleUserFollowers(_ int, _ string) []interface{} {
	return []interface{}{
		Edges{
			URL:        fmt.Sprintf("/users/%d/followers", 3),
			TotalCount: 4,
			NextCursor: "user4",
			Items: []string{
				"5",
			},
		},
		User{
			ID:   "5",
			URL:  "/users/5",
			Name: "user 5",
		},
	}
}

// /users/:id/following
func handleUserFollowing(_ int, _ string) []interface{} {
	return []interface{}{
		Edges{
			URL:        fmt.Sprintf("/users/%d/following", 3),
			TotalCount: 1,
			NextCursor: "user4",
			Items: []string{
				"6",
			},
		},
		User{
			ID:   "6",
			URL:  "/users/6",
			Name: "user 6",
		},
	}
}

// /users/:id/posts
func handleUserPosts(viewerID int, url string) []interface{} {
	userID := strings.TrimPrefix(url, "/users/")
	userID = strings.TrimSuffix(userID, "/posts")

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("feed:%s", userID)
	lenCmd := pipe.LLen(context.Background(), key)
	rangeCmd := pipe.LRange(context.Background(), fmt.Sprintf("feed:%s", userID), 0, 20)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	posts := Edges{
		URL:        fmt.Sprintf("/users/%s/posts", userID),
		TotalCount: int(lenCmd.Val()),
		NextCursor: "",
		Items:      rangeCmd.Val(),
	}

	var results []interface{}
	results = append(results, posts)
	for _, postID := range rangeCmd.Val() {
		results = append(results, handlePostsId(viewerID, postID)...)
	}

	return results
}

// /posts/:id
func handlePostsId(viewerID int, url string) []interface{} {
	postID, _ := strconv.Atoi(strings.TrimPrefix(url, "/posts/"))

	result := Post{
		URL: fmt.Sprintf("/posts/%d", postID),
		ID:  strconv.Itoa(postID),
	}

	post := store.Post{}
	err := store.NodeGet(postID, &post)
	if err != nil {
		result.Text = "DELETED"
		return []interface{}{result}
	}

	result.Text = post.Text
	result.Date = time.Unix(int64(post.Date), 0).Format(time.RFC3339)
	result.UserID = strconv.Itoa(post.UserID)

	user := handleUserById(viewerID, fmt.Sprintf("/users/%d", post.UserID))

	var results []interface{}
	results = append(results, result)
	results = append(results, user...)
	return results
}

// /posts/:id/liked
func handlePostsLiked(_ int, url string) []interface{} {
	postID := strings.TrimPrefix(url, "/posts/")
	postID = strings.TrimSuffix(postID, "/liked")

	key := fmt.Sprintf("postLikes:%s", postID)
	card, _ := store.RedisClient.ZCard(context.Background(), key).Result()

	edge := Edges{
		URL:        fmt.Sprintf("/posts/%s/liked", postID),
		TotalCount: int(card),
		NextCursor: "",
		Items:      []string{},
	}

	var result []interface{}
	result = append(result, edge)

	return result
}

type PostLikeData struct {
	URL        string `json:"url"`
	PostID     string `json:"postId,omitempty"`
	IsLiked    bool   `json:"isLiked,omitempty"`
	LikesCount int    `json:"likesCount,omitempty"`
}

// /posts/:id/isLiked
func handlePostsIsLiked(viewerID int, url string) []interface{} {
	postID := strings.TrimPrefix(url, "/posts/")
	postID = strings.TrimSuffix(postID, "/isLiked")

	key := fmt.Sprintf("postLikes:%s", postID)
	pipe := store.RedisClient.Pipeline()
	cardCmd := pipe.ZCard(context.Background(), key)
	scoreCmd := pipe.ZScore(context.Background(), key, strconv.Itoa(viewerID))

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error getting likes: %s", err)
	}

	edge := PostLikeData{
		URL:        url,
		PostID:     postID,
		IsLiked:    scoreCmd.Val() != 0,
		LikesCount: int(cardCmd.Val()),
	}

	var result []interface{}
	result = append(result, edge)

	return result
}

// /viewer
func handleViewer(viewerID int, _ string) []interface{} {
	type Viewer struct {
		URL      string `json:"url,omitempty"`
		ViewerID string `json:"viewerId,omitempty"`
	}

	viewer := Viewer{
		URL: fmt.Sprintf("/viewer"),
	}
	results := []interface{}{&viewer}

	if viewerID != 0 {
		viewer.ViewerID = strconv.Itoa(viewerID)
		results = append(results, handleUserById(viewerID, fmt.Sprintf("/users/%d", viewerID))...)
	}

	return results
}

func handleQuery(viewerID int, url string) []interface{} {
	type route struct {
		Pattern string
		Handler func(viewerID int, url string) []interface{}
	}

	routes := []route{
		{"/feed", handleFeed},

		{"/users/(\\w+)", handleUserById},
		{"/users/(\\w+)/followers", handleUserFollowers},
		{"/users/(\\w+)/following", handleUserFollowing},
		{"/users/(\\w+)/posts", handleUserPosts},

		{"/posts/(\\w+)", handlePostsId},
		{"/posts/(\\w+)/liked", handlePostsLiked},
		{"/posts/(\\w+)/isLiked", handlePostsIsLiked},

		{"/viewer", handleViewer},
	}

	for _, r := range routes {
		if m, _ := regexp.MatchString("^"+r.Pattern+"$", url); m {
			return r.Handler(viewerID, url)
		}
	}

	return nil
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	authToken := r.Header.Get("authorization")
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	userID, _ := auth.CheckToken(authToken)

	urls := strings.Split(r.URL.Query().Get("urls"), ",")

	results := make([]interface{}, 0)
	for _, query := range urls {
		results = append(results, handleQuery(userID, query)...)
	}

	_ = json.NewEncoder(w).Encode(results)
}
