package api

import (
	"context"
	"net/http"
	"strings"

	authapi "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/auth"
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

func AuthMiddleware(authService *Service, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		authHeader = strings.TrimSpace(authHeader)

		if authHeader != "" {
			verifyReq := &authapi.VerifyTokenRequest{
				Token: authHeader,
			}
			verifyResp, err := authService.VerifyToken(ctx, verifyReq)
			if err == nil {
				ctx = context.WithValue(r.Context(), UserIDKey, verifyResp.UserId)
			}
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

var ErrAuthRequired = twirp.NewError(twirp.Unauthenticated, "auth_required")
