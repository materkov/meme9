package api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/pkg/posts"
	"github.com/materkov/meme9/web5/store"
	"log"
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

// /posts/:id
func handlePostsId(ctx context.Context, viewerID int, url string) []interface{} {
	postID, _ := strconv.Atoi(strings.TrimPrefix(url, "/posts/"))

	result := Post{
		URL: fmt.Sprintf(url),
		ID:  strconv.Itoa(postID),
	}

	post := store.PostStoreFromCtx(ctx).Get(postID)
	if post == nil || !posts.CanSee(post, viewerID) {
		result.IsDeleted = true
		return []interface{}{result}
	}

	result.Text = post.Text
	result.Date = time.Unix(int64(post.Date), 0).UTC().Format(time.RFC3339)
	result.UserID = strconv.Itoa(post.UserID)
	result.CanDelete = post.UserID == viewerID

	var results []interface{}
	results = append(results, result)
	results = append(results, fmt.Sprintf("/users/%d", post.UserID))
	results = append(results, fmt.Sprintf("/posts/%d/liked?count=0", postID))
	return results
}

// /posts/:id/liked
func handlePostsLiked(ctx context.Context, viewerID int, reqURL string) []interface{} {
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

	post := store.PostStoreFromCtx(ctx).Get(postID)
	if !posts.CanSee(post, viewerID) {
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

	_, err := pipe.Exec(context.Background())
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
