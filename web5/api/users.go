package api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/store"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type User struct {
	URL    string `json:"url,omitempty"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Bio    string `json:"bio,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// /users/:id
func handleUserById(ctx context.Context, _ int, url string) []interface{} {
	userID, _ := strconv.Atoi(strings.TrimPrefix(url, "/users/"))

	user := store.CachedStoreFromCtx(ctx).User.Get(userID)

	wrapped := User{
		ID:  strconv.Itoa(userID),
		URL: fmt.Sprintf("/users/%d", userID),
	}

	if user == nil || userID <= 0 {
		wrapped.Name = "Deleted User"
		return []interface{}{
			wrapped,
			fmt.Sprintf("/users/%d/online", userID),
		}
	}

	wrapped.Name = user.Name

	if user.AvatarSha != "" {
		wrapped.Avatar = files.GetURL(user.AvatarSha)
	} else if user.VkPhoto200 != "" {
		wrapped.Avatar = user.VkPhoto200
	}

	return []interface{}{
		wrapped,
		fmt.Sprintf("/users/%d/online", userID),
	}
}

type Online struct {
	URL string `json:"url,omitempty"`

	IsOnline bool `json:"isOnline,omitempty"`
}

// /users/:id/online
func handleUserOnline(ctx context.Context, _ int, url string) []interface{} {
	userID, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(url, "/online"), "/users/"))

	isOnline := store.CachedStoreFromCtx(ctx).Online.Get(userID)

	wrapped := Online{
		URL:      url,
		IsOnline: isOnline,
	}

	return []interface{}{wrapped}
}

// /users/:id/followers
func handleUserFollowers(_ context.Context, viewerID int, url string) []interface{} {
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
func handleUserFollowing(_ context.Context, _ int, url string) []interface{} {
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
func handleUserPosts(_ context.Context, viewerID int, reqURL string) []interface{} {
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
			results = append(results, fmt.Sprintf("/posts/%s", postID))
		}
	}

	return results
}
