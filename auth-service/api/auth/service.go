package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"

	"github.com/materkov/meme9/auth-service/adapters/tokens"
	"github.com/materkov/meme9/auth-service/adapters/users"
	authapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/auth"
	tokensservice "github.com/materkov/meme9/auth-service/services/tokens"
)

type TokensAdapter interface {
	Create(ctx context.Context, token tokens.Token) (string, error)
	GetByValue(ctx context.Context, tokenValue string) (*tokens.Token, error)
}

type UsersAdapter interface {
	GetByUsername(ctx context.Context, username string) (*users.User, error)
	Create(ctx context.Context, user users.User) (string, error)
}

type Service struct {
	users         UsersAdapter
	tokens        TokensAdapter
	tokensService *tokensservice.Service
}

func NewService(usersAdapter UsersAdapter, tokensAdapter TokensAdapter, tokensService *tokensservice.Service) *Service {
	return &Service{
		users:         usersAdapter,
		tokens:        tokensAdapter,
		tokensService: tokensService,
	}
}

func (s *Service) Login(ctx context.Context, req *authapi.LoginRequest) (*authapi.LoginResponse, error) {
	errInvalidCredentials := twirp.NewErrorf(twirp.Unauthenticated, "invalid_credentials")
	if req.Username == "" || req.Password == "" {
		return nil, errInvalidCredentials
	}

	user, err := s.users.GetByUsername(ctx, req.Username)
	if errors.Is(err, users.ErrNotFound) {
		return nil, errInvalidCredentials
	} else if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errInvalidCredentials
	}

	authTokenStr, err := s.createAndSaveAuthToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating auth token: %w", err)
	}

	return &authapi.LoginResponse{
		Token:    authTokenStr,
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

func (s *Service) Register(ctx context.Context, req *authapi.RegisterRequest) (*authapi.LoginResponse, error) {
	if req.Username == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "username_required")
	} else if req.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "password_required")
	}

	user, err := s.users.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, users.ErrNotFound) {
		return nil, fmt.Errorf("error getting user: %w", err)
	} else if user != nil {
		return nil, twirp.NewError(twirp.AlreadyExists, "username_exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error generating password hash: %w", err)
	}

	user = &users.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	userID, err := s.users.Create(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	authTokenStr, err := s.createAndSaveAuthToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error creating auth token: %w", err)
	}

	return &authapi.LoginResponse{
		Token:    authTokenStr,
		UserId:   userID,
		Username: user.Username,
	}, nil
}

func (s *Service) createAndSaveAuthToken(ctx context.Context, userID string) (string, error) {
	authTokenStr, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	authToken := tokens.Token{
		UserID:    userID,
		CreatedAt: time.Now(),
		Token:     authTokenStr,
	}

	_, err = s.tokens.Create(ctx, authToken)
	if err != nil {
		return "", fmt.Errorf("error saving auth token: %w", err)
	}

	return authTokenStr, nil
}

func (s *Service) VerifyToken(ctx context.Context, req *authapi.VerifyTokenRequest) (*authapi.VerifyTokenResponse, error) {
	errInvalidToken := twirp.NewErrorf(twirp.Unauthenticated, "invalid_token")

	tokenValue := strings.TrimPrefix(req.Token, "Bearer ")
	tokenValue = strings.TrimSpace(tokenValue)

	if tokenValue == "" {
		return nil, errInvalidToken
	}

	authToken, err := s.tokens.GetByValue(ctx, tokenValue)
	if errors.Is(err, tokens.ErrNotFound) {
		return nil, errInvalidToken
	} else if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}

	return &authapi.VerifyTokenResponse{
		UserId: authToken.UserID,
	}, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

