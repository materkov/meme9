package apiwrapper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/materkov/meme9/web7/api"
)

// PublishHandler handles post publishing requests
type PublishHandler struct {
	*BaseHandler
}

// NewPublishHandler creates a new publish handler
func NewPublishHandler(api *api.API) *PublishHandler {
	return &PublishHandler{
		BaseHandler: NewBaseHandler(api),
	}
}

type PublishRequest struct {
	Text string `json:"text"`
}

type PublishResponse struct {
	ID string `json:"id"`
}

// Handle processes publish requests
func (h *PublishHandler) Handle(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var reqBody PublishRequest
	if err := json.Unmarshal(body, &reqBody); err != nil {
		writeErrorCode(w, "invalid_json", "")
		return
	}

	userID := GetUserID(req)
	if userID == "" {
		writeErrorCode(w, "unauthorized", "")
		return
	}

	apiReq := api.PublishRequest{
		Text: reqBody.Text,
	}

	resp, err := h.api.Publish(req.Context(), apiReq, userID)
	if err != nil {
		if err.Error() == "text_empty" {
			writeErrorCode(w, "text_empty", "")
			return
		}
		if err.Error() == "text_too_long" {
			writeErrorCode(w, "text_too_long", "")
			return
		}
		if errors.Is(err, errors.New("unauthorized")) {
			writeErrorCode(w, "unauthorized", "")
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	json.NewEncoder(w).Encode(PublishResponse{ID: resp.ID})
}
