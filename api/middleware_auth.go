package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/store"
)

type AuthMiddleware struct {
	store *store.Store
}

func (a *AuthMiddleware) getTokenFromCookie(r *http.Request) (*store.Token, error) {
	tokenCookie, err := r.Cookie("access_token")
	if err != nil {
		return nil, nil
	}

	nodeID, err := store.GetNodeIDFromToken(tokenCookie.Value)
	if err != nil {
		return nil, nil
	}

	token, err := a.store.GetToken(nodeID)
	if err == store.ErrNodeNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error getting token from store: %w", err)
	}

	return token, nil
}

type contextKey string

const viewerCtxKey = contextKey("viewer")

func (a *AuthMiddleware) Do(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := a.getTokenFromCookie(r)
		if err != nil {
			log.Printf("[ERROR] Error checking auth token: %s", err)
		}

		requestScheme := r.Header.Get("x-forwarded-proto")
		if requestScheme == "" {
			requestScheme = "http"
		}

		viewer := api.Viewer{
			UserAgent:     r.Header.Get("user-agent"),
			RequestHost:   r.Host,
			RequestScheme: requestScheme,
		}

		if token != nil {
			viewer.User, err = a.store.GetUser(token.UserID)
			if err != nil {
				log.Printf("[ERROR] Error getting user by id: %s", err)
			}
		}

		newContext := context.WithValue(r.Context(), viewerCtxKey, &viewer)
		next(w, r.WithContext(newContext))
	}
}
