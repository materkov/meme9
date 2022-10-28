package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/pkg/posts"
	"github.com/materkov/meme9/web5/store"
	"log"
	"net/http"
	"net/url"
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

	IsDeleted bool `json:"isDeleted,omitempty"`

	CanDelete bool `json:"canDelete,omitempty"`
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
func handleFeed(viewerID int, reqUrl string) []interface{} {
	parsedURL, _ := url.Parse(reqUrl)
	cursor, _ := strconv.Atoi(parsedURL.Query().Get("cursor"))
	count := 10

	postIds, err := store.RedisClient.LRange(context.Background(), "feed", int64(cursor), int64(cursor+count-1)).Result()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	nextCursor := ""
	if len(postIds) == count {
		nextCursor = strconv.Itoa(cursor + count)
	}

	feed := Edges{
		URL:        reqUrl,
		TotalCount: 20,
		NextCursor: nextCursor,
		Items:      postIds,
	}

	var results []interface{}
	results = append(results, feed)

	for _, postID := range postIds {
		results = append(results, handlePostsId(viewerID, "/posts/"+postID)...)
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
func handleUserFollowers(viewerID int, url string) []interface{} {
	type FollowerEdges struct {
		Edges
		IsFollowing bool `json:"isFollowing,omitempty"`
	}

	pipe := store.RedisClient.Pipeline()

	userID, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(url, "/followers"), "/users/"))
	cardCmd := pipe.ZCard(context.Background(), fmt.Sprintf("followed_by:%d", userID))
	scoreCmd := pipe.ZScore(context.Background(), fmt.Sprintf("followed_by:%d", userID), strconv.Itoa(viewerID))

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error redis: %s", err)
	}

	return []interface{}{
		FollowerEdges{
			Edges: Edges{
				URL:        fmt.Sprintf("/users/%d/followers", userID),
				TotalCount: int(cardCmd.Val()),
				NextCursor: "",
				Items: []string{
					"",
				},
			}, IsFollowing: scoreCmd.Val() != 0,
		},
	}
}

// /users/:id/following
func handleUserFollowing(_ int, url string) []interface{} {
	userID, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(url, "/following"), "/users/"))
	result, _ := store.RedisClient.ZCard(context.Background(), fmt.Sprintf("following:%d", userID)).Result()

	return []interface{}{
		Edges{
			URL:        fmt.Sprintf("/users/%d/following", userID),
			TotalCount: int(result),
			NextCursor: "",
			Items:      []string{},
		},
	}
}

// /users/:id/posts
func handleUserPosts(viewerID int, reqURL string) []interface{} {
	parsedURL, _ := url.Parse(reqURL)
	cursor, _ := strconv.Atoi(parsedURL.Query().Get("cursor"))
	count, _ := strconv.Atoi(parsedURL.Query().Get("count"))

	r := regexp.MustCompile(`^/users/(\w+)/`)
	regexpResults := r.FindStringSubmatch(reqURL)

	userID, _ := strconv.Atoi(regexpResults[1])

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("feed:%d", userID)
	lenCmd := pipe.LLen(context.Background(), key)

	var rangeCmd *redis.StringSliceCmd
	if count > 0 {
		rangeCmd = pipe.LRange(context.Background(), fmt.Sprintf("feed:%d", userID), int64(cursor), int64(cursor+count-1))
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	nextCursor := ""
	if cursor+count < int(lenCmd.Val()) {
		nextCursor = strconv.Itoa(cursor + count)
	}

	edges := Edges{
		URL:        reqURL,
		TotalCount: int(lenCmd.Val()),
		NextCursor: nextCursor,
	}

	if rangeCmd != nil {
		edges.Items = rangeCmd.Val()
	}

	var results []interface{}
	results = append(results, edges)

	if rangeCmd != nil {
		for _, postID := range rangeCmd.Val() {
			results = append(results, handlePostsId(viewerID, postID)...)
		}
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
	if err != nil || !posts.CanSee(&post, viewerID) {
		result.IsDeleted = true
		return []interface{}{result}
	}

	result.Text = post.Text
	result.Date = time.Unix(int64(post.Date), 0).Format(time.RFC3339)
	result.UserID = strconv.Itoa(post.UserID)
	result.CanDelete = post.UserID == viewerID

	user := handleUserById(viewerID, fmt.Sprintf("/users/%d", post.UserID))
	postLiked := handlePostsLiked(viewerID, fmt.Sprintf("/posts/%d/liked?count=0", postID))

	var results []interface{}
	results = append(results, result)
	results = append(results, user...)
	results = append(results, postLiked...)
	return results
}

// /posts/:id/liked
func handlePostsLiked(viewerID int, reqURL string) []interface{} {
	type LikedEdges struct {
		Edges
		IsViewerLiked bool `json:"isViewerLiked,omitempty"`
	}

	r := regexp.MustCompile(`^/posts/(\w+)/`)
	results := r.FindStringSubmatch(reqURL)

	postID, _ := strconv.Atoi(results[1])

	count := 0
	parsedURL, _ := url.Parse(reqURL)
	if parsedURL != nil {
		count, _ = strconv.Atoi(parsedURL.Query().Get("count"))
	}

	edge := LikedEdges{
		Edges: Edges{URL: reqURL},
	}

	post := &store.Post{}
	err := store.NodeGet(postID, post)
	if err != nil || !posts.CanSee(post, viewerID) {
		return []interface{}{edge}
	}

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("postLikes:%d", postID)
	cardCmd := pipe.ZCard(context.Background(), key)
	isLikedCmd := pipe.ZScore(context.Background(), key, strconv.Itoa(viewerID))

	var usersCmd *redis.StringSliceCmd
	if count > 0 {
		usersCmd = pipe.ZRevRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
			Count:  int64(count),
		})
	}

	_, err = pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error redis: %s", err)
	}

	edge.TotalCount = int(cardCmd.Val())
	edge.IsViewerLiked = isLikedCmd.Val() != 0

	if usersCmd != nil {
		edge.Items = usersCmd.Val()
	}

	return []interface{}{edge}
}

// /viewer
func handleViewer(viewerID int, _ string) []interface{} {
	type Viewer struct {
		URL      string `json:"url,omitempty"`
		ViewerID string `json:"viewerId,omitempty"`
	}

	viewer := Viewer{
		URL: "/viewer",
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

		{"/viewer", handleViewer},
	}

	for _, r := range routes {
		path := url
		idx := strings.Index(path, "?")
		if idx != -1 {
			path = path[:idx]
		}

		if m, _ := regexp.MatchString("^"+r.Pattern+"$", path); m {
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
