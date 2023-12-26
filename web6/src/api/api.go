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
	e := json.NewEncoder(w)

	if err != nil {
		errResp := struct {
			Error string `json:"error"`
		}{}

		var publicErr Error
		if ok := errors.As(err, &publicErr); ok {
			errResp.Error = string(publicErr)
		} else {
			errResp.Error = "Internal server error"
		}

		_ = e.Encode(errResp)
	} else {
		_ = e.Encode(resp)
	}
}

type Viewer struct {
	UserID       int
	AuthToken    string
	IsCookieAuth bool
	ClientIP     string
}

type Error string

func (e Error) Error() string {
	return fmt.Sprintf("API Error: %s", string(e))
}
