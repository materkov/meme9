package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type PublishReq struct {
	Text string `json:"text"`
}

type PublishResp struct {
	ID string `json:"id"`
}

func (a *API) publishHandler(w http.ResponseWriter, r *http.Request) {
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

	authHeader := r.Header.Get("Authorization")
	userID, err := a.postsService.VerifyToken(r.Context(), authHeader)
	if err != nil {
		writeUnauthorized(w, "unauthorized")
		return
	}

	post, err := a.postsService.CreatePost(r.Context(), publishReq.Text, userID)
	if err != nil {
		writeInternalServerError(w, "failed to create post")
		return
	}

	json.NewEncoder(w).Encode(PublishResp{ID: post.ID})
}
