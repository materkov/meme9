package api

import (
	"context"
	"net/http"

	"github.com/materkov/meme9/web7/api/auth"
	authapi "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/auth"
	"github.com/twitchtv/twirp"
)

const httpHeadersKey contextKey = "httpHeaders"

// AuthHook creates a Twirp server hook for authentication
func AuthHook(authService *auth.Service) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			// Skip auth for all Auth service methods - they handle their own validation
			// This prevents circular dependency when VerifyToken is called from within the hook
			serviceName, _ := twirp.ServiceName(ctx)
			methodName, _ := twirp.MethodName(ctx)

			if serviceName == "meme.auth.Auth" {
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
				"Publish":     true,
				"Subscribe":   true,
				"Unsubscribe": true,
				"GetStatus":   true,
			}

			requiresAuth := authRequiredMethods[methodName]

			if authHeader == "" {
				if requiresAuth {
					// Return error for methods that require auth
					return ctx, twirp.NewError(twirp.Unauthenticated, "authorization required")
				}
				// Some endpoints don't require auth (like Posts.GetFeed with type="all")
				// We'll let the handler decide
				return ctx, nil
			}

			// Use proto VerifyToken method
			verifyReq := &authapi.VerifyTokenRequest{
				Token: authHeader,
			}
			verifyResp, err := authService.VerifyToken(ctx, verifyReq)
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
			ctx = context.WithValue(ctx, UserIDKey, userID)
			return ctx, nil
		},
	}
}

// CORSMiddleware adds CORS headers to allow requests from frontend
func CORSMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		origin := r.Header.Get("Origin")
		// Allow requests from localhost:3000 (frontend) or any origin in development
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// AuthMiddleware wraps the Twirp handler to inject HTTP headers into context
func AuthMiddleware(handler http.Handler) http.Handler {
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
