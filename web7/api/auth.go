package api

import (
	"context"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

func (a *API) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		userID, err := a.tokensService.VerifyToken(r.Context(), authHeader)
		if err != nil {
			writeErrorCode(w, "unauthorized", "")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func getUserID(r *http.Request) string {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
