// Package main provides JSON utility functions for HTTP request/response handling.
package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// writeJSON writes a JSON response with the given status code and data.
// It sets the appropriate Content-Type header and encodes the data as JSON.
//
// Parameters:
//   - w: HTTP response writer
//   - status: HTTP status code to set
//   - data: Data to encode as JSON
//
// Returns:
//   - error: JSON encoding error if any
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// readJSON reads and decodes JSON data from the request body into the provided data structure.
// It limits the request body size to 1MB and disallows unknown fields for security.
//
// Parameters:
//   - w: HTTP response writer (for setting MaxBytesReader)
//   - r: HTTP request containing JSON data
//   - data: Pointer to struct where decoded JSON will be stored
//
// Returns:
//   - error: JSON decoding error or size limit exceeded error
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 MB limit for request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(data)
}

// writeJSONError writes a JSON error response with the given status code and message.
// It wraps the error message in a standard envelope structure.
//
// Parameters:
//   - w: HTTP response writer
//   - status: HTTP status code to set
//   - message: Error message to include in response
//
// Returns:
//   - error: JSON encoding error if any
func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	log.Printf(message)
	return writeJSON(w, status, &envelope{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}
