package likes

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/api"
	likesapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/likes"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks

type LikesAdapter interface {
	Like(ctx context.Context, postID, userID string) error
	Unlike(ctx context.Context, postID, userID string) error
	IsLiked(ctx context.Context, postID, userID string) (bool, error)
}

type Service struct {
	likes LikesAdapter
}

func NewService(likesAdapter LikesAdapter) *Service {
	return &Service{likes: likesAdapter}
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
