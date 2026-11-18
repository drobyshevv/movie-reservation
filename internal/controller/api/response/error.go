package response

import (
	"encoding/json"
	"net/http"
)

const (
	CodeInvalidRequest  = "invalid_request"
	CodeValidationError = "validation_error"
	CodeNotFound        = "not_found"
	CodeConflict        = "conflict"
	CodeUnauthorized    = "unauthorized"
	CodeInternalError   = "internal_error"
)

type ErrorResponse struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func WriteError(w http.ResponseWriter, typeError, title string, status int, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Type:   typeError,
		Title:  title,
		Status: status,
		Detail: detail,
	})
}
