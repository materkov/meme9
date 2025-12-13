package apiwrapper

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/materkov/meme9/web7/api"
)

// FeedHandler handles feed-related requests
type FeedHandler struct {
	api *api.API
}

// NewFeedHandler creates a new feed handler
func NewFeedHandler(api *api.API) *FeedHandler {
	return &FeedHandler{
		api: api,
	}
}

type FeedPostResponse struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

type FeedRequest struct {
	Type string `json:"type"`
}

// Handle processes feed requests
func (h *FeedHandler) Handle(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	userID, err := h.api.VerifyToken(req.Context(), authHeader)
	if err != nil {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	var reqBody FeedRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	apiReq := api.FeedRequest{
		Type: reqBody.Type,
	}

	feedPosts, err := h.api.GetFeed(req.Context(), apiReq, userID)
	if err != nil {
		if err.Error() == "authentication required for subscriptions feed" {
			writeErrorCode(w, "unauthorized", err.Error())
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	// Convert to response format
	response := make([]FeedPostResponse, len(feedPosts))
	for i, post := range feedPosts {
		response[i] = FeedPostResponse{
			ID:        post.ID,
			Text:      post.Text,
			UserID:    post.UserID,
			Username:  post.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		}
	}

	json.NewEncoder(w).Encode(response)
}
