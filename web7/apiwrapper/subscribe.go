package apiwrapper

import (
	"encoding/json"
	"net/http"

	"github.com/materkov/meme9/web7/api"
)

type SubscribeRequest struct {
	UserID string `json:"user_id"`
}

type SubscribeResponse struct {
	Subscribed bool `json:"subscribed"`
}

func (r *Router) subscribeHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody SubscribeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	followerID := getUserID(req)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	apiReq := api.SubscribeRequest{
		UserID: reqBody.UserID,
	}

	resp, err := r.api.Subscribe(req.Context(), apiReq, followerID)
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

func (r *Router) unsubscribeHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody SubscribeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	followerID := getUserID(req)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	apiReq := api.SubscribeRequest{
		UserID: reqBody.UserID,
	}

	resp, err := r.api.Unsubscribe(req.Context(), apiReq, followerID)
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

func (r *Router) subscriptionStatusHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody SubscribeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	followerID := getUserID(req)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	apiReq := api.SubscribeRequest{
		UserID: reqBody.UserID,
	}

	resp, err := r.api.GetSubscriptionStatus(req.Context(), apiReq, followerID)
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
