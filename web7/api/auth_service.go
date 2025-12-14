package api

import (
	"context"
	"errors"
	"time"

	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"

	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	authapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/auth"
	tokensservice "github.com/materkov/meme9/web7/services/tokens"
)

type AuthService struct {
	users         *users.Adapter
	tokens        *tokens.Adapter
	tokensService *tokensservice.Service
}

func NewAuthService(usersAdapter *users.Adapter, tokensAdapter *tokens.Adapter, tokensService *tokensservice.Service) *AuthService {
	return &AuthService{
		users:         usersAdapter,
		tokens:        tokensAdapter,
		tokensService: tokensService,
	}
}

// Login implements the Auth Login method
func (s *AuthService) Login(ctx context.Context, req *authapi.LoginRequest) (*authapi.LoginResponse, error) {
	if req.Username == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "username_required")
	}
	if req.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "password_required")
	}

	user, err := s.users.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "invalid_credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "invalid_credentials")
	}

	tokenValue, err := generateToken()
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	_, err = s.tokens.Create(ctx, tokens.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	return &authapi.LoginResponse{
		Token:    tokenValue,
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

// Register implements the Auth Register method
func (s *AuthService) Register(ctx context.Context, req *authapi.RegisterRequest) (*authapi.LoginResponse, error) {
	if req.Username == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "username_required")
	}
	if req.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "password_required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	user, err := s.users.Create(ctx, users.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	})
	if err != nil {
		if errors.Is(err, users.ErrUsernameExists) {
			return nil, twirp.NewError(twirp.AlreadyExists, "username_exists")
		}
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	tokenValue, err := generateToken()
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	_, err = s.tokens.Create(ctx, tokens.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, twirp.NewError(twirp.Internal, "internal_server_error")
	}

	return &authapi.LoginResponse{
		Token:    tokenValue,
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

// VerifyToken implements the Auth VerifyToken method
func (s *AuthService) VerifyToken(ctx context.Context, req *authapi.VerifyTokenRequest) (*authapi.VerifyTokenResponse, error) {
	if req.Token == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "token is required")
	}

	userID, err := s.tokensService.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "invalid token")
	}

	return &authapi.VerifyTokenResponse{
		UserId: userID,
	}, nil
}
