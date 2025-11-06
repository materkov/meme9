package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func writeBadRequest(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, message)
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	writeError(w, http.StatusUnauthorized, message)
}

func writeConflict(w http.ResponseWriter, message string) {
	writeError(w, http.StatusConflict, message)
}

func writeInternalServerError(w http.ResponseWriter, message string) {
	writeError(w, http.StatusInternalServerError, message)
}
