package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/materkov/meme9/web7/adapters/mongo"
)

type PublishReq struct {
	Text string `json:"text"`
}

type PublishResp struct {
	ID string `json:"id"`
}

func (a *API) verifyToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	// Support both "Bearer token" and just "token"
	tokenValue := strings.TrimPrefix(authHeader, "Bearer ")
	tokenValue = strings.TrimSpace(tokenValue)

	token, err := a.mongo.GetTokenByValue(r.Context(), tokenValue)
	if err != nil {
		if errors.Is(err, mongo.ErrTokenNotFound) {
			return "", fmt.Errorf("invalid token")
		}
		return "", fmt.Errorf("error verifying token: %w", err)
	}

	return token.UserID, nil
}

func (a *API) publishHandler(w http.ResponseWriter, r *http.Request) {
	// Verify authentication token
	userID, err := a.verifyToken(r)
	if err != nil {
		writeUnauthorized(w, "unauthorized")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	var publishReq PublishReq
	err = json.Unmarshal(body, &publishReq)
	if err != nil {
		writeBadRequest(w, "invalid JSON")
		return
	}

	post, err := a.mongo.AddPost(r.Context(), mongo.Post{
		Text:      publishReq.Text,
		UserID:    userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		writeInternalServerError(w, "failed to create post")
		return
	}

	json.NewEncoder(w).Encode(PublishResp{ID: post.ID})
}
