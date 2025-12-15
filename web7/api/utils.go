package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
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

func RequireError(t *testing.T, err error, code string) {
	t.Helper()

	twirpErr, ok := err.(twirp.Error)
	require.True(t, ok)
	require.Equal(t, code, twirpErr.Msg())
}
