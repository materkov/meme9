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
	URL    string `json:"url,omitempty"`
	ID     string `json:"id,omitempty"`
	Date   string `json:"date,omitempty"`
	Text   string `json:"text,omitempty"`
	UserID string `json:"userId,omitempty"`
	User   *User  `json:"user,omitempty"`

	IsDeleted bool `json:"isDeleted,omitempty"`

	CanDelete bool   `json:"canDelete,omitempty"`
	PhotoID   string `json:"photoId,omitempty"`
	Photo     *Photo `json:"photo,omitempty"`

	LikesConnection *PostsLikesConnection `json:"likesConnection,omitempty"`
}

type PostsList struct {
	Items      []*Post `json:"items,omitempty"`
	NextCursor string  `json:"nextCursor,omitempty"`
}

// /posts/:id
func handlePostsId(ctx context.Context, viewerID int, url string) []interface{} {
	postID, _ := strconv.Atoi(strings.TrimPrefix(url, "/posts/"))

	result := Post{
		URL: url,
		ID:  strconv.Itoa(postID),
	}

	post := store.CachedStoreFromCtx(ctx).Post.Get(postID)
	if post == nil || !posts.CanSee(post, viewerID) {
		result.IsDeleted = true
		return []interface{}{result}
	}

	result.Text = post.Text
	result.Date = time.Unix(int64(post.Date), 0).UTC().Format(time.RFC3339)
	result.UserID = strconv.Itoa(post.UserID)
	result.CanDelete = post.UserID == viewerID

	if post.PhotoID != 0 {
		result.PhotoID = strconv.Itoa(post.PhotoID)
	}

	var results []interface{}
	results = append(results, result)
	results = append(results, fmt.Sprintf("/users/%d", post.UserID))
	results = append(results, fmt.Sprintf("/posts/%d/liked?count=0", postID))

	if post.PhotoID != 0 {
		store.CachedStoreFromCtx(ctx).Photo.Preload([]int{post.PhotoID})
		results = append(results, fmt.Sprintf("/photos/%d", post.PhotoID))
	}

	return results
}

type LikedEdges struct {
	Edges
	IsViewerLiked bool `json:"isViewerLiked,omitempty"`
}

// /posts/:id/liked
func handlePostsLiked(ctx context.Context, viewerID int, reqURL string) []interface{} {

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

	post := store.CachedStoreFromCtx(ctx).Post.Get(postID)
	if !posts.CanSee(post, viewerID) {
		return []interface{}{edge}
	}

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("postLikes:%d", postID)

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

	isLiked, count := store.CachedStoreFromCtx(ctx).Liked.Get(viewerID, postID)

	edge.TotalCount = count
	edge.IsViewerLiked = isLiked

	if usersCmd != nil {
		edge.Items = usersCmd.Val()
	}

	return []interface{}{edge}
}

type PostsAdd struct {
	Text  string `json:"text"`
	Photo string `json:"photo"`
}

func handlePostsAdd(ctx context.Context, viewerID int, req *PostsAdd) (*Post, error) {
	if req.Text == "" {
		return nil, fmt.Errorf("empty text")
	}
	if viewerID == 0 {
		return nil, fmt.Errorf("not authorized")
	}

	photoID, _ := strconv.Atoi(req.Photo)

	postID, err := posts.Add(req.Text, viewerID, photoID)
	if err != nil {
		return nil, fmt.Errorf("error adding post: %w", err)
	}

	result := handlePostsId(ctx, viewerID, fmt.Sprintf("/posts/%d", postID))
	post := result[0].(Post)

	return &post, nil
}

type PostsDelete struct {
	ID string `json:"id"`
}

func handlePostsDelete(ctx context.Context, viewerID int, req *PostsDelete) error {
	postID, _ := strconv.Atoi(req.ID)

	post := &store.Post{}
	err := store.NodeGet(postID, post)
	if err == store.ErrNodeNotFound {
		return fmt.Errorf("post not found")
	} else if err != nil {
		return fmt.Errorf("error getting post: %w", err)
	}

	if post.UserID != viewerID {
		return fmt.Errorf("no access to delete this post")
	}

	if !post.IsDeleted {
		err = posts.Delete(post)
		if err != nil {
			return err
		}
	}

	return nil
}

type PostsLike struct {
	PostID string `json:"postId"`
}

func handlePostsLike(ctx context.Context, viewerID int, req *PostsLike) (*LikedEdges, error) {
	postID, _ := strconv.Atoi(req.PostID)

	post := store.Post{}
	err := store.NodeGet(postID, &post)
	if err == store.ErrNodeNotFound {
		return nil, fmt.Errorf("post not found")
	} else if err != nil {
		return nil, err
	}

	if viewerID == 0 {
		return nil, fmt.Errorf("not authorized")
	}

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("postLikes:%d", postID)
	_ = pipe.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: viewerID,
	})
	cardCmd := pipe.ZCard(context.Background(), key)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return nil, err
	}

	result := LikedEdges{
		Edges: Edges{
			TotalCount: int(cardCmd.Val()),
		},
		IsViewerLiked: true,
	}

	return &result, nil
}

type PostsUnlike struct {
	PostID string `json:"postId"`
}

func handlePostsUnlike(ctx context.Context, viewerID int, req *PostsUnlike) (*LikedEdges, error) {
	postID, _ := strconv.Atoi(req.PostID)

	post := store.Post{}
	err := store.NodeGet(postID, &post)
	if err == store.ErrNodeNotFound {
		return nil, fmt.Errorf("post not found")
	} else if err != nil {
		return nil, err
	}

	if viewerID == 0 {
		return nil, fmt.Errorf("not authorized")
	}

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("postLikes:%d", postID)
	pipe.ZRem(context.Background(), key, viewerID)
	cardCmd := pipe.ZCard(context.Background(), key)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return nil, err
	}

	result := LikedEdges{
		Edges: Edges{
			TotalCount: int(cardCmd.Val()),
		},
		IsViewerLiked: false,
	}

	return &result, nil
}

type PostsGetLikesConnection struct {
	PostID string `json:"postId"`
	Count  int    `json:"count"`
}

type PostsLikesConnection struct {
	TotalCount    int  `json:"totalCount,omitempty"`
	IsViewerLiked bool `json:"isViewerLiked,omitempty"`

	Items []*User `json:"items,omitempty"`
}

func handlePostsLikesConnection(ctx context.Context, viewerID int, req *PostsGetLikesConnection) (*PostsLikesConnection, error) {
	postID, _ := strconv.Atoi(req.PostID)

	post := store.CachedStoreFromCtx(ctx).Post.Get(postID)
	if !posts.CanSee(post, viewerID) {
		return nil, fmt.Errorf("no access to post")
	}

	pipe := store.RedisClient.Pipeline()

	key := fmt.Sprintf("postLikes:%d", postID)

	var usersCmd *redis.StringSliceCmd
	if req.Count > 0 {
		usersCmd = pipe.ZRevRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
			Count:  int64(req.Count),
		})
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Error redis: %s", err)
	}

	isLiked, count := store.CachedStoreFromCtx(ctx).Liked.Get(viewerID, postID)

	result := PostsLikesConnection{
		TotalCount:    count,
		IsViewerLiked: isLiked,
	}

	for _, userID := range usersCmd.Val() {
		users := handleUserById(ctx, viewerID, fmt.Sprintf("/users/%s", userID))
		user := users[0].(User)

		result.Items = append(result.Items, &user)
	}

	return &result, nil
}
