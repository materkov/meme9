package api

import (
	"context"
	"fmt"

	"github.com/twitchtv/twirp"

	usersapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/users"
	"github.com/materkov/meme9/users-service/adapters/users"
)

type UsersAdapter interface {
	GetByID(ctx context.Context, userID string) (*users.User, error)
	UpdateAvatar(ctx context.Context, userID, avatarURL string) error
}

type Service struct {
	users UsersAdapter
}

func NewService(usersAdapter UsersAdapter) *Service {
	return &Service{
		users: usersAdapter,
	}
}

func (s *Service) Get(ctx context.Context, req *usersapi.GetUserRequest) (*usersapi.GetUserResponse, error) {
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_not_found")
	}

	user, err := s.users.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, "user_not_found")
	}

	return &usersapi.GetUserResponse{
		Id:        user.ID,
		Username:  user.Username,
		AvatarUrl: user.AvatarURL,
	}, nil
}

func (s *Service) SetAvatar(ctx context.Context, req *usersapi.SetAvatarRequest) (*usersapi.SetAvatarResponse, error) {
	// Check authentication
	authenticatedUserID := GetUserIDFromContext(ctx)
	if authenticatedUserID == "" {
		return nil, ErrAuthRequired
	}

	// Validate request
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id_required")
	}

	// Users can only set their own avatar
	if req.UserId != authenticatedUserID {
		return nil, twirp.NewError(twirp.PermissionDenied, "can_only_set_own_avatar")
	}

	// Validate avatar URL
	if req.AvatarUrl == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "avatar_url_required")
	}

	// Update avatar
	err := s.users.UpdateAvatar(ctx, req.UserId, req.AvatarUrl)
	if err != nil {
		if err == users.ErrNotFound {
			return nil, twirp.NewError(twirp.NotFound, "user_not_found")
		}
		return nil, twirp.NewError(twirp.Internal, fmt.Sprintf("failed to update avatar: %v", err))
	}

	return &usersapi.SetAvatarResponse{}, nil
}
