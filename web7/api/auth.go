package api

import (
	"context"
	"net/http"

	"github.com/materkov/meme9/web7/api/auth"
	authapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/auth"
	"github.com/twitchtv/twirp"
)

func AuthHook(authService *auth.Service) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			header, ok := twirp.HTTPRequestHeaders(ctx)
			authHeader := ""
			if ok && header != nil {
				authHeader = header.Get("Authorization")
			}

			if authHeader != "" {
				verifyReq := &authapi.VerifyTokenRequest{
					Token: authHeader,
				}
				verifyResp, err := authService.VerifyToken(ctx, verifyReq)
				if err == nil {
					userID := verifyResp.UserId
					ctx = context.WithValue(ctx, UserIDKey, userID)
				}
			}

			return ctx, nil
		},
	}
}

func CORSMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

var ErrAuthRequired = twirp.NewError(twirp.Unauthenticated, "auth_required")
