package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// config holds all configuration parameters for the application
type config struct {
	// addr is the network address the server will listen on
	addr string
	// shutdownTimeout defines the maximum time allowed for graceful shutdown
	shutdownTimeout time.Duration
}

// application holds the dependencies for HTTP handlers, helpers, and middleware.
// It serves as the dependency injection container for the entire application.
type application struct {
	config config
	logger *slog.Logger
}

// mount configures and returns the HTTP router with all middleware and routes.
// It sets up a production-ready middleware stack including request ID, logging,
// recovery, real IP detection, and request timeouts.
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Production-ready middleware stack
	r.Use(middleware.RequestID) // Adds unique request ID for tracing
	r.Use(middleware.RealIP)    // Sets RemoteAddr to real client IP
	r.Use(middleware.Logger)    // Logs request details
	r.Use(middleware.Recoverer) // Recovers from panics and returns 500

	r.Use(middleware.Timeout(60 * time.Second))

	// API versioning with grouped routes
	r.Route("/v1", func(r chi.Router) {
		// Health check endpoints
		r.Get("/health", app.healthcheckHandler)
	})

	return r
}

// getLocalIP returns the first non-loopback IPv4 address of the local machine.
// This is used for debugging and logging purposes to identify which instance
// is serving requests in multi-instance deployments.
func (app *application) getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		app.logger.Warn("Failed to get local IP addresses", "error", err)
		return "unknown"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	app.logger.Warn("No non-loopback IPv4 address found")
	return "unknown"
}
