package twirp

import (
	"context"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/materkov/meme9/web7/api"
	json_api "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/json_api"
)

type contextKey string

const httpHeadersKey contextKey = "httpHeaders"

// AuthHook creates a Twirp server hook for authentication
func AuthHook(apiAdapter *api.API) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			// Skip auth for login and register
			serviceName, _ := twirp.ServiceName(ctx)
			methodName, _ := twirp.MethodName(ctx)

			if serviceName == "meme.json_api.JsonAPI" &&
				(methodName == "Login" || methodName == "Register") {
				return ctx, nil
			}

			// Extract token from header
			// Try Twirp's function first, then fall back to our custom context value
			header, ok := twirp.HTTPRequestHeaders(ctx)
			if !ok || header == nil {
				// Fallback to our custom context value
				if headers, ok := ctx.Value(httpHeadersKey).(http.Header); ok {
					header = headers
				}
			}

			authHeader := ""
			if header != nil {
				authHeader = header.Get("Authorization")
			}

			// Methods that require authentication
			authRequiredMethods := map[string]bool{
				"Publish":               true,
				"Subscribe":             true,
				"Unsubscribe":           true,
				"GetSubscriptionStatus": true,
			}

			requiresAuth := authRequiredMethods[methodName]

			if authHeader == "" {
				if requiresAuth {
					// Return error for methods that require auth
					return ctx, twirp.NewError(twirp.Unauthenticated, "authorization required")
				}
				// Some endpoints don't require auth (like GetFeed with type="all")
				// We'll let the handler decide
				return ctx, nil
			}

			// Use proto VerifyToken method
			verifyReq := &json_api.VerifyTokenRequest{
				Token: authHeader,
			}
			verifyResp, err := apiAdapter.VerifyToken(ctx, verifyReq)
			if err != nil {
				if requiresAuth {
					// Return error for methods that require auth
					return ctx, twirp.NewError(twirp.Unauthenticated, "invalid token")
				}
				// For optional auth endpoints, let the handler decide
				return ctx, nil
			}
			userID := verifyResp.UserId

			// Add user ID to context using api package's context key
			ctx = context.WithValue(ctx, api.UserIDKey, userID)
			return ctx, nil
		},
	}
}

// AuthMiddleware wraps the Twirp handler to inject HTTP headers into context
func AuthMiddleware(apiAdapter *api.API, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// Store headers in context using our custom key
		ctx = context.WithValue(ctx, httpHeadersKey, r.Header)

		// Also try to set them using Twirp's function if available
		if newCtx, err := twirp.WithHTTPRequestHeaders(ctx, r.Header); err == nil && newCtx != nil {
			ctx = newCtx
		}

		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}
