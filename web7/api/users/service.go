package users

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/users"
	usersapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/users"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks

type UsersAdapter interface {
	GetByID(ctx context.Context, userID string) (*users.User, error)
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
		Id:       user.ID,
		Username: user.Username,
	}, nil
}
