package posts

import (
	"context"
	"fmt"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
)

//go:generate mockgen -source=posts.go -destination=mocks/posts_adapter_mock.go -package=mocks

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
