package api

import (
	"context"
	"net/http"

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
		ctx := context.WithValue(r.Context(), UserIDKey, r.Header.Get("x-user-id"))
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

var ErrAuthRequired = twirp.NewError(twirp.Unauthenticated, "auth_required")
