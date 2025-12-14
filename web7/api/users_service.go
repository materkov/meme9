package api

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/adapters/users"
	usersapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/users"
)

type UsersService struct {
	users *users.Adapter
}

func NewUsersService(usersAdapter *users.Adapter) *UsersService {
	return &UsersService{
		users: usersAdapter,
	}
}

// Get implements the Users Get method
func (s *UsersService) Get(ctx context.Context, req *usersapi.GetUserRequest) (*usersapi.GetUserResponse, error) {
	if req.UserId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	user, err := s.users.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, "user not found")
	}

	return &usersapi.GetUserResponse{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}
