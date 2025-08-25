package main

import (
	"net/http"
	"time"
)

// HealthResponse represents the health check response structure
//
//	@Description	Health status information for the microservice
type HealthResponse struct {
	// Current operational status of the service
	//	@example	"healthy"
	Status string `json:"status" example:"healthy"`

	// Descriptive message about the service state
	//	@example	"Service is operational"
	Message string `json:"message" example:"Service is operational"`

	// UTC timestamp when the health check was performed
	//	@example	"2024-01-15T10:30:00.123Z"
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00.123Z"`

	// Service version information (optional)
	//	@example	"1.0.0"
	Version string `json:"version,omitempty" example:"1.0.0"`
}

// HealthCheck returns the current health status of the service
//
//	@Summary		Get service health status
//	@Description	Returns comprehensive health information about the microservice including operational status,
//	@Description	system timestamp, and version details. This endpoint is used for monitoring and load balancer health checks.
//	@Description	Always returns 200 OK when the service is running and can process requests.
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthResponse		"Service is healthy and operational"
//	@Failure		500	{object}	map[string]string	"Internal server error - service may be unhealthy"
//	@Router			/v1/health [get]
//	@x-order		1
//
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

	if err := app.jsonResponse(w, http.StatusOK, healthResp); err != nil {
		app.logger.Error("Failed to encode health check response", "error", err)
		app.internalServerError(w, r, err)
		return
	}

	app.logger.Info("Health check completed successfully",
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent())
}
