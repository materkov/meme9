package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/twitchtv/twirp"
)

type contextKey string

const UserIDKey contextKey = "userID"

func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user ID from x-user-id header (set by frontend proxy after token verification)
		userID := r.Header.Get("x-user-id")
		userID = strings.TrimSpace(userID)

		if userID != "" {
			ctx = context.WithValue(r.Context(), UserIDKey, userID)
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

var ErrAuthRequired = twirp.NewError(twirp.Unauthenticated, "auth_required")
