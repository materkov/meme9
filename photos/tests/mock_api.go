package tests

import (
	"context"
	"fmt"

	authpb "github.com/materkov/meme9/photos/internal/authclient/pb/github.com/materkov/meme9/api/auth"
	twirp "github.com/twitchtv/twirp"
)

type MockAuth struct{}

func (m *MockAuth) Login(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return nil, fmt.Errorf("MockAuth.Login: unimplemented")
}

func (m *MockAuth) Register(context.Context, *authpb.RegisterRequest) (*authpb.LoginResponse, error) {
	return nil, fmt.Errorf("MockAuth.Register: unimplemented")
}

func (m *MockAuth) VerifyToken(_ context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	if req.Token == "good-token" {
		return &authpb.VerifyTokenResponse{
			UserId: "test-user",
		}, nil
	}

	return nil, twirp.NewError(twirp.Unauthenticated, "invalid_token")
}
