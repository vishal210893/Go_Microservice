package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// ===== API configuration =====
type api struct {
	addr string // Server address in format "host:port"
}

// ===== Domain model =====
type User struct {
	Name  string `json:"name"`  // User's full name
	Age   int    `json:"age"`   // User's age
	Email string `json:"email"` // User's email address
}

// ===== In-memory user store =====
// Using a mutex to keep it safe when multiple requests happen at the same time
var (
	users []User
	mu    sync.RWMutex
)

// ===== Handlers =====

// getUsersHandler returns all users as JSON
func (s *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Lock for safe concurrent read
	mu.RLock()
	defer mu.RUnlock()

	// Always set status before writing body
	w.WriteHeader(http.StatusOK)

	// Encode users slice into JSON
	_ = json.NewEncoder(w).Encode(users)
}

// createUsersHandler adds a new user from request body
func (s *api) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Prevent unknown JSON fields

	// Parse JSON body into User struct
	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Add new user safely
	mu.Lock()
	users = append(users, payload)
	mu.Unlock()

	// Return created user in response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(payload)
}

// ===== main =====

func main() {
	api := &api{addr: ":8080"} // Initialize API with server address

	// Create a new HTTP router (Go 1.22+ patterns)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUsersHandler)

	// Start server
	fmt.Printf("Server starting on %s\n", api.addr)
	if err := http.ListenAndServe(api.addr, mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
