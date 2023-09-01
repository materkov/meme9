package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type API struct{}

type Void struct{}

var ErrParsingRequest = &Error{
	Code:    400,
	Message: "failed to parse request",
}

func writeResp(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)

		var publicErr *Error
		if ok := errors.As(err, &publicErr); ok {
			fmt.Fprint(w, publicErr.Message)
		} else {
			fmt.Fprint(w, "Internal server error")
		}
	} else {
		_ = json.NewEncoder(w).Encode(resp)
	}
}
