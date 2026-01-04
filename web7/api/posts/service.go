package posts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	postsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/posts"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks

type PostsAdapter interface {
	Add(ctx context.Context, post posts.Post) (*posts.Post, error)
	GetByUserID(ctx context.Context, userID string) ([]posts.Post, error)
	GetByID(ctx context.Context, postID string) (*posts.Post, error)
	GetAll(ctx context.Context) ([]posts.Post, error)
	GetByUserIDs(ctx context.Context, userIDs []string) ([]posts.Post, error)
}

type UsersAdapter interface {
	GetByID(ctx context.Context, userID string) (*users.User, error)
	GetByIDs(ctx context.Context, userIDs []string) (map[string]*users.User, error)
}

type SubscriptionsAdapter interface {
	GetFollowing(ctx context.Context, userID string) ([]string, error)
}

type LikesAdapter interface {
	GetLikesCounts(ctx context.Context, postIDs []string) (map[string]int32, error)
	GetLikedByUser(ctx context.Context, userID string, postIDs []string) (map[string]bool, error)
}

type Service struct {
	posts         PostsAdapter
	users         UsersAdapter
	subscriptions SubscriptionsAdapter
	likes         LikesAdapter
}

func NewService(postsAdapter PostsAdapter, usersAdapter UsersAdapter, subscriptions SubscriptionsAdapter, likesAdapter LikesAdapter) *Service {
	return &Service{
		posts:         postsAdapter,
		users:         usersAdapter,
		subscriptions: subscriptions,
		likes:         likesAdapter,
	}
}

func (s *Service) Publish(ctx context.Context, req *postsapi.PublishRequest) (*postsapi.PublishResponse, error) {
	userID := api.GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, api.ErrAuthRequired
	}

	if req.Text == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "text_empty")
	}
	if len(req.Text) > 1000 {
		return nil, twirp.NewError(twirp.InvalidArgument, "text_too_long")
	}

	post, err := s.posts.Add(ctx, posts.Post{
		Text:      req.Text,
		UserID:    userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &postsapi.PublishResponse{
		Id: post.ID,
	}, nil
}

func (s *Service) GetByUsers(ctx context.Context, req *postsapi.GetByUsersRequest) (*postsapi.GetByUsersResponse, error) {
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id_required")
	}

	postsListChan := make(chan []posts.Post)
	go func() {
		postsList, err := s.posts.GetByUserID(ctx, req.UserId)
		if err != nil {
			log.Printf("failed to get posts: %s", err)
			postsListChan <- []posts.Post{}
		} else {
			postsListChan <- postsList
		}
	}()

	usernameChan := make(chan string)
	go func() {
		user, err := s.users.GetByID(ctx, req.UserId)
		if err == nil && user != nil {
			usernameChan <- user.Username
		} else {
			log.Printf("failed to get user: %s", err)
			usernameChan <- ""
		}
	}()

	postsList := <-postsListChan
	username := <-usernameChan

	userPosts := make([]*postsapi.Post, len(postsList))
	postIDs := make([]string, len(postsList))
	for i, post := range postsList {
		postIDs[i] = post.ID
	}

	userID := api.GetUserIDFromContext(ctx)
	likesCounts := make(map[string]int32)
	likedByUser := make(map[string]bool)
	counts, err := s.likes.GetLikesCounts(ctx, postIDs)
	if err == nil {
		likesCounts = counts
	}
	if userID != "" {
		liked, err := s.likes.GetLikedByUser(ctx, userID, postIDs)
		if err == nil {
			likedByUser = liked
		}
	}

	for i, post := range postsList {
		userPosts[i] = makeProtoPost(&post, username, likesCounts[post.ID], likedByUser[post.ID])
	}

	return &postsapi.GetByUsersResponse{
		Posts: userPosts,
	}, nil
}

func (s *Service) Get(ctx context.Context, req *postsapi.GetPostRequest) (*postsapi.Post, error) {
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id_required")
	}

	post, err := s.posts.GetByID(ctx, req.PostId)
	if errors.Is(err, posts.ErrNotFound) {
		return nil, twirp.NewError(twirp.NotFound, "post_not_found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	userName := ""
	user, err := s.users.GetByID(ctx, post.UserID)
	if err != nil {
		log.Printf("Cannot load user info: %s", err)
	} else {
		userName = user.Username
	}

	userID := api.GetUserIDFromContext(ctx)
	var likesCount int32
	var isLiked bool
	count, err := s.likes.GetLikesCounts(ctx, []string{post.ID})
	if err == nil {
		likesCount = count[post.ID]
	}
	if userID != "" {
		liked, err := s.likes.GetLikedByUser(ctx, userID, []string{post.ID})
		if err == nil {
			isLiked = liked[post.ID]
		}
	}

	return makeProtoPost(post, userName, likesCount, isLiked), nil
}

func (s *Service) GetFeed(ctx context.Context, req *postsapi.FeedRequest) (*postsapi.FeedResponse, error) {
	userID := api.GetUserIDFromContext(ctx)

	feedType := req.Type
	if feedType == postsapi.FeedType_FEED_TYPE_UNSPECIFIED {
		feedType = postsapi.FeedType_FEED_TYPE_ALL
	}

	var postsList []posts.Post
	var err error

	if feedType == postsapi.FeedType_FEED_TYPE_SUBSCRIPTIONS {
		if userID == "" {
			return nil, api.ErrAuthRequired
		}

		followingIDs, err := s.subscriptions.GetFollowing(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get following ids: %w", err)
		}

		subscribedUserIDs := append(followingIDs, userID)
		postsList, err = s.posts.GetByUserIDs(ctx, subscribedUserIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get posts: %w", err)
		}
	} else {
		postsList, err = s.posts.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get posts: %w", err)
		}
	}

	userIDSet := make(map[string]bool)
	for _, post := range postsList {
		if post.UserID != "" {
			userIDSet[post.UserID] = true
		}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	usersMap, err := s.users.GetByIDs(ctx, userIDs)
	if err != nil {
		log.Printf("failed to get users: %s", err)
		usersMap = make(map[string]*users.User)
	}

	postIDs := make([]string, len(postsList))
	for i, post := range postsList {
		postIDs[i] = post.ID
	}

	likesCounts := make(map[string]int32)
	likedByUser := make(map[string]bool)
	counts, err := s.likes.GetLikesCounts(ctx, postIDs)
	if err == nil {
		likesCounts = counts
	}
	if userID != "" {
		liked, err := s.likes.GetLikedByUser(ctx, userID, postIDs)
		if err == nil {
			likedByUser = liked
		}
	}

	feedPosts := make([]*postsapi.Post, len(postsList))
	for i, post := range postsList {
		userName := ""
		if user := usersMap[post.UserID]; user != nil {
			userName = user.Username
		}

		feedPosts[i] = makeProtoPost(&post, userName, likesCounts[post.ID], likedByUser[post.ID])
	}

	return &postsapi.FeedResponse{
		Posts: feedPosts,
	}, nil
}

func makeProtoPost(post *posts.Post, userName string, likesCount int32, isLiked bool) *postsapi.Post {
	return &postsapi.Post{
		Id:         post.ID,
		Text:       post.Text,
		UserId:     post.UserID,
		UserName:   userName,
		CreatedAt:  post.CreatedAt.Format(time.RFC3339),
		LikesCount: likesCount,
		IsLiked:    isLiked,
	}
}
