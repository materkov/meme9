package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/tokens"
)

type PostsService struct {
	postsAdapter  *posts.Adapter
	tokensAdapter *tokens.Adapter
}

func NewPostsService(postsAdapter *posts.Adapter, tokensAdapter *tokens.Adapter) *PostsService {
	return &PostsService{
		postsAdapter:  postsAdapter,
		tokensAdapter: tokensAdapter,
	}
}

func (s *PostsService) VerifyToken(ctx context.Context, authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	// Support both "Bearer token" and just "token"
	tokenValue := strings.TrimPrefix(authHeader, "Bearer ")
	tokenValue = strings.TrimSpace(tokenValue)

	token, err := s.tokensAdapter.GetByValue(ctx, tokenValue)
	if err != nil {
		if errors.Is(err, tokens.ErrNotFound) {
			return "", fmt.Errorf("invalid token")
		}
		return "", fmt.Errorf("error verifying token: %w", err)
	}

	return token.UserID, nil
}

func (s *PostsService) CreatePost(ctx context.Context, text string, userID string) (*posts.Post, error) {
	post, err := s.postsAdapter.Add(ctx, posts.Post{
		Text:      text,
		UserID:    userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}
