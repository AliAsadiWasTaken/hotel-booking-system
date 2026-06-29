package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// validate is a package-level singleton to avoid the cost of
// rebuilding the validator's reflection cache on every request.
var validate = validator.New()

type errorResponse struct {
	Error string `json:"error"`
}

// WriteJSON encodes data as JSON and writes it to the response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, errorResponse{Error: message})
}

// Decode reads and decodes a JSON request body into v.
func Decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// Validate runs struct-tag validation on v and returns a human-readable
// message for the first failing field, or an empty string if all pass.
func Validate(v any) string {
	err := validate.Struct(v)
	if err == nil {
		return ""
	}

	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			return e.Field() + " is required"
		case "min":
			return e.Field() + " must be at least " + e.Param()
		case "email":
			return e.Field() + " must be a valid email address"
		case "gt":
			return e.Field() + " must be greater than " + e.Param()
		}
	}

	return err.Error()
}
