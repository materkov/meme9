package apiwrapper

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/materkov/meme9/web7/api"
)

// UserPostsHandler handles user posts requests
type UserPostsHandler struct {
	api *api.API
}

// NewUserPostsHandler creates a new user posts handler
func NewUserPostsHandler(api *api.API) *UserPostsHandler {
	return &UserPostsHandler{
		api: api,
	}
}

type UserPostsRequest struct {
	UserID string `json:"user_id"`
}

type UserPostResponse struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

// Handle processes user posts requests
func (h *UserPostsHandler) Handle(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var reqBody UserPostsRequest
	if err := json.Unmarshal(body, &reqBody); err != nil {
		writeErrorCode(w, "invalid_json", "")
		return
	}

	apiReq := api.UserPostsRequest{
		UserID: reqBody.UserID,
	}

	userPosts, err := h.api.GetUserPosts(req.Context(), apiReq)
	if err != nil {
		if err.Error() == "user_id is required" {
			writeErrorCode(w, "user_id_required", "")
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	response := make([]UserPostResponse, len(userPosts))
	for i, post := range userPosts {
		response[i] = UserPostResponse{
			ID:        post.ID,
			Text:      post.Text,
			UserID:    post.UserID,
			Username:  post.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	json.NewEncoder(w).Encode(response)
}
