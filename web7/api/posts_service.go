package api

import (
	"context"
	"errors"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/users"
	postsapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/posts"
	postsservice "github.com/materkov/meme9/web7/services/posts"
)

type PostsService struct {
	posts        *posts.Adapter
	users        *users.Adapter
	postsService *postsservice.Service
}

func NewPostsService(postsAdapter *posts.Adapter, usersAdapter *users.Adapter, postsService *postsservice.Service) *PostsService {
	return &PostsService{
		posts:        postsAdapter,
		users:        usersAdapter,
		postsService: postsService,
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

	userPosts := make([]*postsapi.UserPostResponse, len(postsList))
	for i, post := range postsList {
		userPosts[i] = &postsapi.UserPostResponse{
			Id:        post.ID,
			Text:      post.Text,
			UserId:    post.UserID,
			Username:  username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	return &postsapi.GetByUsersResponse{
		Posts: userPosts,
	}, nil
}

// Get implements the Posts Get method
func (s *PostsService) Get(ctx context.Context, req *postsapi.GetPostRequest) (*postsapi.GetPostResponse, error) {
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id is required")
	}

	post, err := s.posts.GetByID(ctx, req.PostId)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, "post not found")
	}

	return &postsapi.GetPostResponse{
		Id:        post.ID,
		Text:      post.Text,
		UserId:    post.UserID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}, nil
}
