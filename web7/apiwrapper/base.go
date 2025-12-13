package apiwrapper

import (
	"context"
	"net/http"

	"github.com/materkov/meme9/web7/api"
)

// BaseHandler provides common functionality for all route handlers
type BaseHandler struct {
	api *api.API
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(api *api.API) *BaseHandler {
	return &BaseHandler{api: api}
}

type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware creates middleware that verifies authentication token
func (h *BaseHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		userID, err := h.api.VerifyToken(req.Context(), authHeader)
		if err != nil {
			writeErrorCode(w, "unauthorized", "")
			return
		}

		ctx := context.WithValue(req.Context(), userIDKey, userID)
		next(w, req.WithContext(ctx))
	}
}

// GetUserID extracts user ID from request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// JSONMiddleware sets JSON content type header
func JSONMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
