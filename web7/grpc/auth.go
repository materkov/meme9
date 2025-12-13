package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/materkov/meme9/web7/api"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthInterceptor creates a gRPC interceptor for authentication
func AuthInterceptor(api *api.API) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for login and register
		if info.FullMethod == "/meme.json_api.JsonAPI/Login" || info.FullMethod == "/meme.json_api.JsonAPI/Register" {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			// Some endpoints don't require auth (like GetFeed with type="all")
			// We'll let the handler decide
			return handler(ctx, req)
		}

		token := authHeaders[0]
		userID, err := api.VerifyToken(ctx, token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		// Add user ID to context
		ctx = context.WithValue(ctx, userIDKey, userID)
		return handler(ctx, req)
	}
}

// getUserIDFromContext extracts user ID from context
func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
