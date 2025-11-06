package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	postsservice "github.com/materkov/meme9/web7/services/posts"
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
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var publishReq PublishReq
	err = json.Unmarshal(body, &publishReq)
	if err != nil {
		writeErrorCode(w, "invalid_json", "")
		return
	}

	userID := getUserID(r)
	post, err := a.postsService.CreatePost(r.Context(), publishReq.Text, userID)
	if err != nil {
		// Handle validation errors with 400 Bad Request
		if errors.Is(err, postsservice.ErrTextEmpty) {
			writeErrorCode(w, "text_empty", "")
			return
		}
		if errors.Is(err, postsservice.ErrTextTooLong) {
			writeErrorCode(w, "text_too_long", "")
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(PublishResp{ID: post.ID})
}
