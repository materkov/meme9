package apiwrapper

import (
	"context"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

func (r *Router) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		userID, err := r.api.VerifyToken(req.Context(), authHeader)
		if err != nil {
			writeErrorCode(w, "unauthorized", "")
			return
		}

		ctx := context.WithValue(req.Context(), userIDKey, userID)
		next(w, req.WithContext(ctx))
	}
}

func getUserID(r *http.Request) string {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
