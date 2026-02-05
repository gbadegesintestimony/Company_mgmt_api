package utils

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, err APIError) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}
