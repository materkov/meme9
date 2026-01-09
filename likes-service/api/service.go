package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/twitchtv/twirp"

	likesapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/likes"
	usersapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/users"
)

type LikesAdapter interface {
	Like(ctx context.Context, postID, userID string) error
	Unlike(ctx context.Context, postID, userID string) error
	IsLiked(ctx context.Context, postID, userID string) (bool, error)
	GetLikers(ctx context.Context, postID, pageToken string, count int) ([]string, string, error)
	GetLikedByUser(ctx context.Context, userID string, postIDs []string) (map[string]bool, error)
	GetLikesCounts(ctx context.Context, postIDs []string) (map[string]int32, error)
}

type Service struct {
	likes LikesAdapter
}

func NewService(likesAdapter LikesAdapter) *Service {
	return &Service{
		likes: likesAdapter,
	}
}

func (s *Service) Like(ctx context.Context, req *likesapi.LikeRequest) (*likesapi.LikeResponse, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrAuthRequired
	}
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id_required")
	}

	err := s.likes.Like(ctx, req.PostId, userID)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &likesapi.LikeResponse{
		Liked: true,
	}, nil
}

func (s *Service) Unlike(ctx context.Context, req *likesapi.LikeRequest) (*likesapi.LikeResponse, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrAuthRequired
	}
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id_required")
	}

	err := s.likes.Unlike(ctx, req.PostId, userID)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &likesapi.LikeResponse{
		Liked: false,
	}, nil
}

func (s *Service) GetLikers(ctx context.Context, req *likesapi.GetLikersRequest) (*likesapi.GetLikersResponse, error) {
	if req.PostId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_id_required")
	}

	count := req.Count
	if count <= 0 {
		count = 20
	} else if count > 100 {
		count = 100
	}

	userIDs, nextPageToken, err := s.likes.GetLikers(ctx, req.PostId, req.PageToken, int(count))
	if err != nil {
		return nil, fmt.Errorf("failed to get likers: %w", err)
	}

	// Call users service to get user details
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		usersServiceURL = "http://localhost:8082"
	}

	usersClient := usersapi.NewUsersProtobufClient(usersServiceURL, &http.Client{})
	usersMap := make(map[string]*usersapi.GetUserResponse)
	for _, userID := range userIDs {
		userReq := &usersapi.GetUserRequest{UserId: userID}
		userResp, err := usersClient.Get(ctx, userReq)
		if err != nil {
			log.Printf("failed to get user %s: %s", userID, err)
			continue
		}
		usersMap[userID] = userResp
	}

	// Build response
	likers := make([]*likesapi.GetLikersResponse_Liker, 0, len(userIDs))
	for _, userID := range userIDs {
		username := ""
		userAvatar := ""
		if user := usersMap[userID]; user != nil {
			username = user.Username
			userAvatar = user.AvatarUrl
		}

		likers = append(likers, &likesapi.GetLikersResponse_Liker{
			UserId:     userID,
			Username:   username,
			UserAvatar: userAvatar,
		})
	}

	return &likesapi.GetLikersResponse{
		Likers:    likers,
		PageToken: nextPageToken,
	}, nil
}

func (s *Service) CheckLike(ctx context.Context, req *likesapi.CheckLikeRequest) (*likesapi.CheckLikeResponse, error) {
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id_required")
	}
	if len(req.PostIds) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "post_ids_required")
	}

	// Get liked status for the user
	likedMap, err := s.likes.GetLikedByUser(ctx, req.UserId, req.PostIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked status: %w", err)
	}

	// Get likes counts for all posts
	likesCountsMap, err := s.likes.GetLikesCounts(ctx, req.PostIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get likes counts: %w", err)
	}

	// Build response arrays in the same order as postIds
	liked := make([]bool, len(req.PostIds))
	likesCount := make([]int32, len(req.PostIds))
	for i, postID := range req.PostIds {
		liked[i] = likedMap[postID]
		likesCount[i] = likesCountsMap[postID]
	}

	return &likesapi.CheckLikeResponse{
		Liked:      liked,
		LikesCount: likesCount,
	}, nil
}
