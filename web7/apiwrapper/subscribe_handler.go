package apiwrapper

import (
	"encoding/json"
	"net/http"

	"github.com/materkov/meme9/web7/api"
)

// SubscribeHandler handles subscription-related requests
type SubscribeHandler struct {
	*BaseHandler
}

// NewSubscribeHandler creates a new subscribe handler
func NewSubscribeHandler(api *api.API) *SubscribeHandler {
	return &SubscribeHandler{
		BaseHandler: NewBaseHandler(api),
	}
}

type SubscribeRequest struct {
	UserID string `json:"user_id"`
}

type SubscribeResponse struct {
	Subscribed bool `json:"subscribed"`
}

// HandleSubscribe processes subscribe requests
func (h *SubscribeHandler) HandleSubscribe(w http.ResponseWriter, req *http.Request) {
	var reqBody SubscribeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	followerID := GetUserID(req)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	apiReq := api.SubscribeRequest{
		UserID: reqBody.UserID,
	}

	resp, err := h.api.Subscribe(req.Context(), apiReq, followerID)
	if err != nil {
		if err.Error() == "user_id is required" {
			writeErrorCode(w, "invalid_request", "user_id is required")
			return
		}
		if err.Error() == "unauthorized" {
			writeErrorCode(w, "unauthorized", "")
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(SubscribeResponse{Subscribed: resp.Subscribed})
}

// HandleUnsubscribe processes unsubscribe requests
func (h *SubscribeHandler) HandleUnsubscribe(w http.ResponseWriter, req *http.Request) {
	var reqBody SubscribeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	followerID := GetUserID(req)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	apiReq := api.SubscribeRequest{
		UserID: reqBody.UserID,
	}

	resp, err := h.api.Unsubscribe(req.Context(), apiReq, followerID)
	if err != nil {
		if err.Error() == "user_id is required" {
			writeErrorCode(w, "invalid_request", "user_id is required")
			return
		}
		if err.Error() == "unauthorized" {
			writeErrorCode(w, "unauthorized", "")
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(SubscribeResponse{Subscribed: resp.Subscribed})
}

// HandleSubscriptionStatus processes subscription status requests
func (h *SubscribeHandler) HandleSubscriptionStatus(w http.ResponseWriter, req *http.Request) {
	var reqBody SubscribeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	followerID := GetUserID(req)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	apiReq := api.SubscribeRequest{
		UserID: reqBody.UserID,
	}

	resp, err := h.api.GetSubscriptionStatus(req.Context(), apiReq, followerID)
	if err != nil {
		if err.Error() == "user_id is required" {
			writeErrorCode(w, "invalid_request", "user_id is required")
			return
		}
		if err.Error() == "unauthorized" {
			writeErrorCode(w, "unauthorized", "")
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(SubscribeResponse{Subscribed: resp.Subscribed})
}
