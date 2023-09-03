package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type API struct{}

type Void struct{}

var ErrParsingRequest = Error("failed to parse request")

func writeResp(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)

		var publicErr Error
		if ok := errors.As(err, &publicErr); ok {
			fmt.Fprint(w, string(publicErr))
		} else {
			fmt.Fprint(w, "Internal server error")
		}
	} else {
		_ = json.NewEncoder(w).Encode(resp)
	}
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
