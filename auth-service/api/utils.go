package api

import (
	"context"

	"github.com/twitchtv/twirp"
)

type contextKey string

const UserIDKey contextKey = "userID"

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

var ErrAuthRequired = twirp.NewError(twirp.Unauthenticated, "auth_required")
