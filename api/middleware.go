package api

import (
	"context"
	"fmt"
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

func (a *AuthMiddleware) Do(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := a.getTokenFromCookie(r)
		if err != nil {
			panic(1)
		}

		viewer := api.Viewer{}
		if token != nil {
			viewer.UserID = token.UserID
		}

		newContext := context.WithValue(r.Context(), "viewer", &viewer)
		next(w, r.WithContext(newContext))
	}
}
