package posts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
)

//go:generate mockgen -source=posts.go -destination=mocks/posts_adapter_mock.go -package=mocks

const maxTextLength = 1000

var (
	ErrTextEmpty   = errors.New("text cannot be empty")
	ErrTextTooLong = errors.New("text cannot be longer than 1000 characters")
)

type PostsAdapter interface {
	Add(ctx context.Context, post posts.Post) (*posts.Post, error)
}

type Service struct {
	postsAdapter PostsAdapter
}

func New(postsAdapter PostsAdapter) *Service {
	return &Service{
		postsAdapter: postsAdapter,
	}
}

func (s *Service) CreatePost(ctx context.Context, text string, userID string) (*posts.Post, error) {
	if err := validateText(text); err != nil {
		return nil, err
	}

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

func validateText(text string) error {
	if text == "" {
		return ErrTextEmpty
	}
	if len(text) > maxTextLength {
		return ErrTextTooLong
	}
	return nil
}
