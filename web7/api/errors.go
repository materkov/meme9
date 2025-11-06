package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error        string `json:"error"`
	ErrorDetails string `json:"error_details"`
}

func writeError(w http.ResponseWriter, statusCode int, errorCode string, errorDetails string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:        errorCode,
		ErrorDetails: errorDetails,
	})
}

func writeErrorCode(w http.ResponseWriter, errorCode string, errorDetails string) {
	writeError(w, http.StatusBadRequest, errorCode, errorDetails)
}

func writeInternalServerError(w http.ResponseWriter, errorCode string, errorDetails string) {
	writeError(w, http.StatusInternalServerError, errorCode, errorDetails)
}
