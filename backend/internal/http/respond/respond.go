package respond

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

type ErrorEnvelope struct {
	Error APIError `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Fail(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorEnvelope{Error: APIError{Code: code, Message: message}})
}