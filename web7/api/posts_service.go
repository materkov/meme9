package api

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/adapters/users"
	postsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/posts"
	postsservice "github.com/materkov/meme9/web7/services/posts"
)

type PostsService struct {
	posts         *posts.Adapter
	users         *users.Adapter
	postsService  *postsservice.Service
	subscriptions *subscriptions.Adapter
}

func NewPostsService(postsAdapter *posts.Adapter, usersAdapter *users.Adapter, postsService *postsservice.Service, subscriptionsAdapter *subscriptions.Adapter) *PostsService {
	return &PostsService{
		posts:         postsAdapter,
		users:         usersAdapter,
		postsService:  postsService,
		subscriptions: subscriptionsAdapter,
	}
}

// Publish implements the Posts Publish method
func (s *PostsService) Publish(ctx context.Context, req *postsapi.PublishRequest) (*postsapi.PublishResponse, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, twirp.NewError(twirp.Unauthenticated, "unauthorized")
	}

	post, err := s.postsService.CreatePost(ctx, req.Text, userID)
	if err != nil {
		if errors.Is(err, postsservice.ErrTextEmpty) {
			return nil, twirp.NewError(twirp.InvalidArgument, "text_empty")
		}
		if errors.Is(err, postsservice.ErrTextTooLong) {
			return nil, twirp.NewError(twirp.InvalidArgument, "text_too_long")
		}
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &postsapi.PublishResponse{
		Id: post.ID,
	}, nil
}

// GetByUsers implements the Posts GetByUsers method
func (s *PostsService) GetByUsers(ctx context.Context, req *postsapi.GetByUsersRequest) (*postsapi.GetByUsersResponse, error) {
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id_required")
	}

	postsList, err := s.posts.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	user, err := s.users.GetByID(ctx, req.UserId)
	username := ""
	if err == nil && user != nil {
		username = user.Username
	}

	userPosts := make([]*postsapi.Post, len(postsList))
	for i, post := range postsList {
		userPosts[i] = &postsapi.Post{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			UserName:  username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &postsapi.GetByUsersResponse{
		Posts: userPosts,
	}, nil
}

// Get implements the Posts Get method
func (s *PostsService) Get(ctx context.Context, req *postsapi.GetPostRequest) (*postsapi.Post, error) {
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id is required")
	}

	post, err := s.posts.GetByID(ctx, req.PostId)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, "post not found")
	}

	userName := ""
	user, err := s.users.GetByID(ctx, post.UserID)
	if err != nil {
		log.Printf("Cannot load user info: %s", err)
	} else {
		userName = user.Username
	}

	return &postsapi.Post{
		Id:        post.ID,
		Text:      post.Text,
		UserId:    post.UserID,
		UserName:  userName,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetFeed implements the Posts GetFeed method
func (s *PostsService) GetFeed(ctx context.Context, req *postsapi.FeedRequest) (*postsapi.FeedResponse, error) {
	userID := GetUserIDFromContext(ctx)

	feedType := req.Type
	if feedType == "" {
		feedType = "all"
	}

	var postsList []posts.Post
	var err error

	if feedType == "subscriptions" {
		if userID == "" {
			return nil, twirp.NewError(twirp.Unauthenticated, "authentication required for subscriptions feed")
		}

		followingIDs, err := s.subscriptions.GetFollowing(ctx, userID)
		if err != nil {
			followingIDs = []string{}
		}

		subscribedUserIDs := append(followingIDs, userID)
		postsList, err = s.posts.GetByUserIDs(ctx, subscribedUserIDs)
		if err != nil {
			return nil, twirp.NewError(twirp.Internal, err.Error())
		}
	} else {
		postsList, err = s.posts.GetAll(ctx)
		if err != nil {
			return nil, twirp.NewError(twirp.Internal, err.Error())
		}
	}

	// Collect unique user IDs
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
		usersMap = make(map[string]*users.User)
	}

	// Build feed posts using Post type
	feedPosts := make([]*postsapi.Post, len(postsList))
	for i, post := range postsList {
		userName := ""
		if user := usersMap[post.UserID]; user != nil {
			userName = user.Username
		}
		feedPosts[i] = &postsapi.Post{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			UserName:  userName,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &postsapi.FeedResponse{
		Posts: feedPosts,
	}, nil
}
