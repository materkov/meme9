package api

import (
	"encoding/json"
	"net/http"
)

type SubscribeRequest struct {
	UserID string `json:"user_id"`
}

type SubscribeResponse struct {
	Subscribed bool `json:"subscribed"`
}

func (a *API) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorCode(w, "method_not_allowed", "")
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	if req.UserID == "" {
		writeErrorCode(w, "invalid_request", "user_id is required")
		return
	}

	followerID := getUserID(r)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	err := a.subscriptions.Subscribe(r.Context(), followerID, req.UserID)
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(SubscribeResponse{Subscribed: true})
}

func (a *API) unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorCode(w, "method_not_allowed", "")
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	if req.UserID == "" {
		writeErrorCode(w, "invalid_request", "user_id is required")
		return
	}

	followerID := getUserID(r)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	err := a.subscriptions.Unsubscribe(r.Context(), followerID, req.UserID)
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(SubscribeResponse{Subscribed: false})
}

func (a *API) subscriptionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorCode(w, "method_not_allowed", "")
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorCode(w, "invalid_request", "")
		return
	}

	if req.UserID == "" {
		writeErrorCode(w, "invalid_request", "user_id is required")
		return
	}

	userID := req.UserID

	followerID := getUserID(r)
	if followerID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	isSubscribed, err := a.subscriptions.IsSubscribed(r.Context(), followerID, userID)
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(SubscribeResponse{Subscribed: isSubscribed})
}
