package api

//go:generate mockgen -destination=mocks/mock_posts_adapter.go -package=mocks github.com/materkov/meme9/posts-service/api PostsAdapter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/twitchtv/twirp"

	likesapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/likes"
	postsapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/posts"
	usersapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/users"
	"github.com/materkov/meme9/posts-service/adapters/posts"
)

type PostsAdapter interface {
	Add(ctx context.Context, post posts.Post) (*posts.Post, error)
	GetByID(ctx context.Context, postID string) (*posts.Post, error)
	GetAll(ctx context.Context) ([]posts.Post, error)
	GetByUserIDs(ctx context.Context, userIDs []string) ([]posts.Post, error)
	MarkAsDeleted(ctx context.Context, postID string) error
}

type Service struct {
	posts PostsAdapter
}

func NewService(postsAdapter PostsAdapter) *Service {
	return &Service{
		posts: postsAdapter,
	}
}

func (s *Service) getUsersServiceClient() usersapi.Users {
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		usersServiceURL = "http://localhost:8082"
	}
	return usersapi.NewUsersProtobufClient(usersServiceURL, &http.Client{})
}

func (s *Service) getLikesServiceClient() likesapi.Likes {
	likesServiceURL := os.Getenv("LIKES_SERVICE_URL")
	if likesServiceURL == "" {
		likesServiceURL = "http://localhost:8084"
	}
	return likesapi.NewLikesProtobufClient(likesServiceURL, &http.Client{})
}

func (s *Service) Publish(ctx context.Context, req *postsapi.PublishRequest) (*postsapi.PublishResponse, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrAuthRequired
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

	postsList, err := s.posts.GetByUserIDs(ctx, []string{req.UserId})
	if err != nil {
		log.Printf("failed to get posts: %s", err)
		postsList = []posts.Post{}
	}

	// Get user info
	usersClient := s.getUsersServiceClient()
	userResp, err := usersClient.Get(ctx, &usersapi.GetUserRequest{UserId: req.UserId})
	username := ""
	userAvatar := ""
	if err == nil && userResp != nil {
		username = userResp.Username
		userAvatar = userResp.AvatarUrl
	}

	userPosts := make([]*postsapi.Post, len(postsList))
	postIDs := make([]string, len(postsList))
	for i, post := range postsList {
		postIDs[i] = post.ID
	}

	likesCounts := make(map[string]int32)
	likedByUser := make(map[string]bool)

	// Get likes counts and liked status using CheckLike
	if len(postIDs) > 0 {
		likesClient := s.getLikesServiceClient()
		// Get current user ID if authenticated
		currentUserID := GetUserIDFromContext(ctx)

		checkLikeReq := &likesapi.CheckLikeRequest{
			UserId:  currentUserID,
			PostIds: postIDs,
		}
		checkLikeResp, err := likesClient.CheckLike(ctx, checkLikeReq)
		if err == nil && checkLikeResp != nil {
			// Map responses back to post IDs
			for i, postID := range postIDs {
				if i < len(checkLikeResp.Liked) {
					likedByUser[postID] = checkLikeResp.Liked[i]
				}
				if i < len(checkLikeResp.LikesCount) {
					likesCounts[postID] = checkLikeResp.LikesCount[i]
				}
			}
		} else {
			// On error, initialize with defaults
			log.Printf("failed to get likes: %s", err)
			for _, postID := range postIDs {
				likesCounts[postID] = 0
				likedByUser[postID] = false
			}
		}
	}

	for i, post := range postsList {
		userPosts[i] = makeProtoPost(&post, username, userAvatar, likesCounts[post.ID], likedByUser[post.ID])
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
	} else if post.Deleted {
		return nil, twirp.NewError(twirp.NotFound, "post_not_found")
	}

	// Get user info
	usersClient := s.getUsersServiceClient()
	userResp, err := usersClient.Get(ctx, &usersapi.GetUserRequest{UserId: post.UserID})
	userName := ""
	userAvatar := ""
	if err == nil && userResp != nil {
		userName = userResp.Username
		userAvatar = userResp.AvatarUrl
	}

	var likesCount int32
	var isLiked bool

	// Get likes count and liked status using CheckLike
	currentUserID := GetUserIDFromContext(ctx)

	likesClient := s.getLikesServiceClient()
	checkLikeReq := &likesapi.CheckLikeRequest{
		UserId:  currentUserID,
		PostIds: []string{post.ID},
	}
	checkLikeResp, err := likesClient.CheckLike(ctx, checkLikeReq)
	if err == nil && checkLikeResp != nil && len(checkLikeResp.Liked) > 0 {
		isLiked = checkLikeResp.Liked[0]
		if len(checkLikeResp.LikesCount) > 0 {
			likesCount = checkLikeResp.LikesCount[0]
		}
	} else {
		// On error, set defaults
		log.Printf("failed to get likes: %s", err)
		likesCount = 0
		isLiked = false
	}

	return makeProtoPost(post, userName, userAvatar, likesCount, isLiked), nil
}

func (s *Service) GetFeed(ctx context.Context, req *postsapi.FeedRequest) (*postsapi.FeedResponse, error) {
	userID := GetUserIDFromContext(ctx)

	feedType := req.Type
	if feedType == postsapi.FeedType_FEED_TYPE_UNSPECIFIED {
		feedType = postsapi.FeedType_FEED_TYPE_ALL
	}

	var postsList []posts.Post
	var err error

	if feedType == postsapi.FeedType_FEED_TYPE_SUBSCRIPTIONS {
		if userID == "" {
			return nil, ErrAuthRequired
		}

		// Get following IDs from subscriptions service
		// Note: We need a GetFollowing method in subscriptions service
		// For now, we'll use GetAll
		postsList, err = s.posts.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get posts: %w", err)
		}
	} else {
		postsList, err = s.posts.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get posts: %w", err)
		}
	}

	// Get unique user IDs
	userIDSet := make(map[string]bool)
	for _, post := range postsList {
		if post.UserID != "" {
			userIDSet[post.UserID] = true
		}
	}

	// Get user info for all users
	usersClient := s.getUsersServiceClient()
	usersMap := make(map[string]*usersapi.GetUserResponse)
	for userID := range userIDSet {
		userResp, err := usersClient.Get(ctx, &usersapi.GetUserRequest{UserId: userID})
		if err == nil && userResp != nil {
			usersMap[userID] = userResp
		}
	}

	postIDs := make([]string, len(postsList))
	for i, post := range postsList {
		postIDs[i] = post.ID
	}

	likesCounts := make(map[string]int32)
	likedByUser := make(map[string]bool)

	// Get likes counts and liked status using CheckLike
	if len(postIDs) > 0 {
		likesClient := s.getLikesServiceClient()
		// Get current user ID if authenticated (userID is already available from context check above)
		currentUserID := userID

		checkLikeReq := &likesapi.CheckLikeRequest{
			UserId:  currentUserID,
			PostIds: postIDs,
		}
		checkLikeResp, err := likesClient.CheckLike(ctx, checkLikeReq)
		if err == nil && checkLikeResp != nil {
			// Map responses back to post IDs
			for i, postID := range postIDs {
				if i < len(checkLikeResp.Liked) {
					likedByUser[postID] = checkLikeResp.Liked[i]
				}
				if i < len(checkLikeResp.LikesCount) {
					likesCounts[postID] = checkLikeResp.LikesCount[i]
				}
			}
		} else {
			// On error, initialize with defaults
			log.Printf("failed to get likes: %s", err)
			for _, postID := range postIDs {
				likesCounts[postID] = 0
				likedByUser[postID] = false
			}
		}
	}

	feedPosts := make([]*postsapi.Post, len(postsList))
	for i, post := range postsList {
		userName := ""
		userAvatar := ""
		if user := usersMap[post.UserID]; user != nil {
			userName = user.Username
			userAvatar = user.AvatarUrl
		}

		feedPosts[i] = makeProtoPost(&post, userName, userAvatar, likesCounts[post.ID], likedByUser[post.ID])
	}

	return &postsapi.FeedResponse{
		Posts: feedPosts,
	}, nil
}

func (s *Service) Delete(ctx context.Context, req *postsapi.DeleteRequest) (*postsapi.DeleteResponse, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrAuthRequired
	}
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id_required")
	}

	post, err := s.posts.GetByID(ctx, req.PostId)
	if errors.Is(err, posts.ErrNotFound) {
		return nil, twirp.NewError(twirp.NotFound, "post_not_found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post.UserID != userID {
		return nil, twirp.NewError(twirp.PermissionDenied, "not_post_owner")
	}

	err = s.posts.MarkAsDeleted(ctx, req.PostId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete post: %w", err)
	}

	return &postsapi.DeleteResponse{}, nil
}

func makeProtoPost(post *posts.Post, userName string, userAvatar string, likesCount int32, isLiked bool) *postsapi.Post {
	return &postsapi.Post{
		Id:         post.ID,
		Text:       post.Text,
		UserId:     post.UserID,
		UserName:   userName,
		UserAvatar: userAvatar,
		CreatedAt:  post.CreatedAt.Format(time.RFC3339),
		LikesCount: likesCount,
		IsLiked:    isLiked,
	}
}
