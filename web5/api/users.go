package api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/imgproxy"
	"github.com/materkov/meme9/web5/pkg/users"
	"github.com/materkov/meme9/web5/store"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type User struct {
	URL    string `json:"url,omitempty"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Bio    string `json:"bio,omitempty"`
	Avatar string `json:"avatar,omitempty"`

	Online *Online `json:"online,omitempty"`
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
		wrapped.Avatar = imgproxy.GetURL(user.AvatarSha, 200)
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

type UsersFollow struct {
	UserID string `json:"userId"`
}

func handleUsersFollow(ctx context.Context, viewerID int, req *UsersFollow) error {
	userID, _ := strconv.Atoi(req.UserID)

	if viewerID == 0 {
		return fmt.Errorf("not authorized")
	} else if userID == 0 {
		return fmt.Errorf("empty user")
	} else if userID == viewerID {
		return fmt.Errorf("you cannot subscribe to yourself")
	}

	err := users.Follow(viewerID, userID)
	if err != nil {
		return err
	}

	return nil
}

type UsersUnfollow struct {
	UserID string `json:"userId"`
}

func handleUsersUnfollow(ctx context.Context, viewerID int, req *UsersUnfollow) error {
	userID, _ := strconv.Atoi(req.UserID)

	err := users.Unfollow(viewerID, userID)
	if err != nil {
		return err
	}

	return nil
}

type UsersEdit struct {
	UserID string
	Name   string
}

func handleUsersEdit(ctx context.Context, viewerID int, req *UsersEdit) error {
	userID, _ := strconv.Atoi(req.UserID)

	user := store.User{}
	err := store.NodeGet(userID, &user)
	if err != nil {
		return err
	} else if user.ID == 0 {
		return fmt.Errorf("user not found")
	}

	if viewerID != user.ID {
		return fmt.Errorf("no access to edit this user")
	}

	if req.Name == "" {
		return fmt.Errorf("name is empty")
	} else if len(req.Name) > 100 {
		return fmt.Errorf("name is too long")
	}

	user.Name = req.Name

	err = store.NodeSave(user.ID, user)
	if err != nil {
		return err
	}

	return nil
}

type UsersSetOnline struct {
}

func handleUsersSetOnline(ctx context.Context, viewerID int, req *UsersSetOnline) error {
	if viewerID == 0 {
		return nil
	}

	go func() {
		_, err := store.RedisClient.Set(context.Background(), fmt.Sprintf("online:%d", viewerID), time.Now().Unix(), time.Minute*3).Result()
		if err != nil {
			log.Printf("Err: %s", err)
		}
	}()

	return nil
}

type UsersSetAvatar struct {
	UploadToken string `json:"uploadToken"`
}

func handleUsersSetAvatar(ctx context.Context, viewerID int, req *UsersSetAvatar) (*User, error) {
	if viewerID == 0 {
		return nil, fmt.Errorf("not authorized")
	}

	user := &store.User{}
	err := store.NodeGet(viewerID, user)
	if err != nil {
		return nil, fmt.Errorf("error getting user")
	}

	photoID, _ := strconv.Atoi(req.UploadToken)
	photo := &store.Photo{}
	err = store.NodeGet(photoID, photo)
	if err != nil {
		return nil, fmt.Errorf("error getting photo")
	}

	user.AvatarSha = photo.Hash

	err = store.NodeSave(user.ID, user)
	if err != nil {
		return nil, fmt.Errorf("error saving user")
	}

	result := handleUserById(ctx, user.ID, fmt.Sprintf("/users/%d", user.ID))
	userWrapped := result[0].(User)

	return &userWrapped, nil
}

type UsersPostsList struct {
	UserID string `json:"userId,omitempty"`
}

func handleUsersPostsList(ctx context.Context, viewerID int, req *UsersPostsList) (*PostsList, error) {
	result := handleUserPosts(ctx, viewerID, fmt.Sprintf("/users/%s/posts?count=10", req.UserID))
	edges := result[0].(Edges)

	response := PostsList{}
	response.Items = wrapPostsList(ctx, viewerID, edges.Items)

	return &response, nil
}
