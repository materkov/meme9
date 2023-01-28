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
		ID: strconv.Itoa(userID),
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
	IsOnline bool `json:"isOnline,omitempty"`
}

// /users/:id/online
func handleUserOnline(ctx context.Context, _ int, url string) []interface{} {
	userID, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(url, "/online"), "/users/"))

	isOnline := store.CachedStoreFromCtx(ctx).Online.Get(userID)

	wrapped := Online{
		IsOnline: isOnline,
	}

	return []interface{}{wrapped}
}

type UserFollowers struct {
	UserID string `json:"userId"`
}

type FollowerEdges struct {
	Edges
	IsFollowing bool `json:"isFollowing,omitempty"`
}

// /users/:id/followers
func handleUserFollowers(_ context.Context, viewerID int, req *UserFollowers) (*FollowerEdges, error) {
	pipe := store.RedisClient.Pipeline()

	//userID, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(url, "/followers"), "/users/"))
	cardCmd := pipe.ZCard(context.Background(), fmt.Sprintf("followed_by:%s", req.UserID))
	scoreCmd := pipe.ZScore(context.Background(), fmt.Sprintf("followed_by:%s", req.UserID), strconv.Itoa(viewerID))

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error redis: %s", err)
	}

	return &FollowerEdges{
		Edges: Edges{
			TotalCount: int(cardCmd.Val()),
			NextCursor: "",
			Items: []string{
				"",
			},
		}, IsFollowing: scoreCmd.Val() != 0,
	}, nil
}

type UserFollowing struct {
	UserID string `json:"userId"`
}

// /users/:id/following
func handleUserFollowing(_ context.Context, _ int, req *UserFollowing) (*Edges, error) {
	result, _ := store.RedisClient.ZCard(context.Background(), fmt.Sprintf("following:%s", req.UserID)).Result()

	return &Edges{
		TotalCount: int(result),
		NextCursor: "",
		Items:      []string{},
	}, nil
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

	user := store.CachedStoreFromCtx(ctx).User.Get(userID)
	if user == nil {
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

	err := store.NodeUpdate(user.ID, user)
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

	user := store.CachedStoreFromCtx(ctx).User.Get(viewerID)
	if user == nil {
		return nil, fmt.Errorf("error getting user")
	}

	photoID, _ := strconv.Atoi(req.UploadToken)
	photo := store.CachedStoreFromCtx(ctx).Photo.Get(photoID)
	if photo == nil {
		return nil, fmt.Errorf("error getting photo")
	}

	user.AvatarSha = photo.Hash

	err := store.NodeUpdate(user.ID, user)
	if err != nil {
		return nil, fmt.Errorf("error saving user")
	}

	result := handleUserById(ctx, user.ID, fmt.Sprintf("/users/%d", user.ID))
	userWrapped := result[0].(User)

	return &userWrapped, nil
}

type UsersPostsList struct {
	UserID string `json:"userId,omitempty"`
	Count  int    `json:"count"`
}

func handleUsersPostsList(ctx context.Context, viewerID int, req *UsersPostsList) (*PostsList, error) {
	result := handleUserPosts(ctx, viewerID, fmt.Sprintf("/users/%s/posts?count=%d", req.UserID, req.Count))
	edges := result[0].(Edges)

	response := PostsList{}
	response.Items = wrapPostsList(ctx, viewerID, edges.Items)
	response.TotalCount = edges.TotalCount

	return &response, nil
}

type UsersList struct {
	UserIds []string `json:"userIds"`
}

func handleUsersList(ctx context.Context, viewerID int, req *UsersList) ([]*User, error) {
	var ids []int
	for _, idStr := range req.UserIds {
		id, _ := strconv.Atoi(idStr)
		if id > 0 {
			ids = append(ids, id)
		}
	}

	store.CachedStoreFromCtx(ctx).User.Preload(ids)

	var results []*User
	for _, userID := range ids {
		user := store.CachedStoreFromCtx(ctx).User.Get(userID)
		if user == nil {
			continue
		}

		xresults := handleUserById(ctx, viewerID, fmt.Sprintf("/users/%d", userID))
		wrappedUser := xresults[0].(User)
		results = append(results, &wrappedUser)
	}

	return results, nil
}
