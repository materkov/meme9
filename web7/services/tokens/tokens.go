package tokens

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/materkov/meme9/web7/adapters/tokens"
)

type Service struct {
	tokensAdapter *tokens.Adapter
}

func New(tokensAdapter *tokens.Adapter) *Service {
	return &Service{
		tokensAdapter: tokensAdapter,
	}
}

func (s *Service) VerifyToken(ctx context.Context, authHeader string) (string, error) {
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
