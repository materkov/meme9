package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type UserPostsReq struct {
	UserID string `json:"user_id"`
}

func (a *API) userPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorCode(w, "method_not_allowed", "")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var req UserPostsReq
	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorCode(w, "invalid_json", "")
		return
	}

	if req.UserID == "" {
		writeErrorCode(w, "user_id_required", "")
		return
	}

	postsList, err := a.posts.GetByUserID(r.Context(), req.UserID)
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	// Fetch user info
	user, err := a.users.GetByID(r.Context(), req.UserID)
	if err != nil {
		// Log error but continue with empty username
		log.Printf("Error fetching user: %v", err)
	}

	username := ""
	if user != nil {
		username = user.Username
	}

	// Map posts to API posts with username
	apiPosts := make([]Post, len(postsList))
	for i, post := range postsList {
		apiPosts[i] = mapPostToAPIPost(post, username)
	}

	json.NewEncoder(w).Encode(apiPosts)
}
