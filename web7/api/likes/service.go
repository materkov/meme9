package likes

import (
	"context"
	"fmt"
	"log"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	likesapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/likes"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks

type LikesAdapter interface {
	Like(ctx context.Context, postID, userID string) error
	Unlike(ctx context.Context, postID, userID string) error
	IsLiked(ctx context.Context, postID, userID string) (bool, error)
	GetLikers(ctx context.Context, postID, pageToken string, count int) ([]string, string, error)
}

type UsersAdapter interface {
	GetByIDs(ctx context.Context, userIDs []string) (map[string]*users.User, error)
}

type Service struct {
	likes LikesAdapter
	users UsersAdapter
}

func NewService(likesAdapter LikesAdapter, usersAdapter UsersAdapter) *Service {
	return &Service{
		likes: likesAdapter,
		users: usersAdapter,
	}
}

func (s *Service) Like(ctx context.Context, req *likesapi.LikeRequest) (*likesapi.LikeResponse, error) {
	userID := api.GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, api.ErrAuthRequired
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
	userID := api.GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, api.ErrAuthRequired
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

	usersMap, err := s.users.GetByIDs(ctx, userIDs)
	if err != nil {
		log.Printf("failed to get users: %s", err)
	}

	// Build response
	likers := make([]*likesapi.GetLikersResponse_Liker, 0, len(userIDs))
	for _, userID := range userIDs {
		username := ""
		if user := usersMap[userID]; user != nil {
			username = user.Username
		}

		likers = append(likers, &likesapi.GetLikersResponse_Liker{
			UserId:   userID,
			Username: username,
		})
	}

	return &likesapi.GetLikersResponse{
		Likers:    likers,
		PageToken: nextPageToken,
	}, nil
}
