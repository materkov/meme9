package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type API struct{}

type Void struct{}

var ErrParsingRequest = Error("FailedParsingRequest")

func writeResp(w http.ResponseWriter, resp interface{}, err error) {
	httpResp := struct {
		Data  interface{} `json:"data"`
		Error string      `json:"error,omitempty"`
	}{}

	if err != nil {
		var publicErr Error
		if ok := errors.As(err, &publicErr); ok {
			httpResp.Error = string(publicErr)
		} else {
			httpResp.Error = "Internal server error"
		}
	} else {
		httpResp.Data = resp
	}

	_ = json.NewEncoder(w).Encode(httpResp)
}

type Viewer struct {
	UserID    int
	AuthToken string
}

type Error string

func (e Error) Error() string {
	return fmt.Sprintf("API Error: %s", string(e))
}

type ctxKey string

const (
	ctxViewer ctxKey = "viewer"
)
