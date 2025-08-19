package main

import (
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// healthcheckHandler handles GET /v1/health requests and returns the service health status.
// It follows industry standards for health check endpoints by returning structured JSON
// with proper HTTP status codes and consistent response format.
//
// Returns:
//   - 200 OK with health information when service is healthy
//   - 500 Internal Server Error if unable to generate response
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	healthResp := HealthResponse{
		Status:    "healthy",
		Message:   "Service is operational",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
	}

	w.WriteHeader(http.StatusOK)

	if err := writeJSON(w, http.StatusOK, healthResp); err != nil {
		app.logger.Error("Failed to encode health check response", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to get health check response")
		return
	}

	app.logger.Info("Health check completed successfully",
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent())
}
