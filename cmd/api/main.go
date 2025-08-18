package main

import (
	"Go-Microservice/internal"
	"Go-Microservice/internal/env"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// DefaultServerTimeout defines the maximum time allowed for graceful server shutdown
	DefaultServerTimeout = 30 * time.Second
	DefaultAddr          = ":8080"
)

// main is the application entry point that initializes the server configuration,
// sets up structured logging, mounts routes, and starts the HTTP server with
// graceful shutdown capabilities for production deployment.
func main() {

	logger := slog.New(internal.NewFormattedLogHandler(os.Stdout, slog.LevelInfo))
	slog.SetDefault(logger)

	app := &application{
		config: config{
			addr:            env.GetPort("ADDR", DefaultAddr),
			shutdownTimeout: DefaultServerTimeout,
		},
		logger: logger,
	}

	router := app.mount()

	if err := app.runWithGracefulShutdown(router); err != nil {
		logger.Error("Failed to start application", "error", err)
		os.Exit(1)
	}
}

// runWithGracefulShutdown starts the HTTP server and implements graceful shutdown
// on receiving termination signals (SIGINT, SIGTERM). This ensures ongoing requests
// are completed before server termination, preventing data loss or corruption.
func (app *application) runWithGracefulShutdown(handler http.Handler) error {
	// Configure HTTP server with production-ready timeouts
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      handler,
		WriteTimeout: 15 * time.Second, // Response write timeout
		ReadTimeout:  10 * time.Second, // Request read timeout
		IdleTimeout:  time.Minute,      // Keep-alive timeout
	}

	// Channel to signal shutdown completion
	shutdownComplete := make(chan struct{})

	// Start shutdown handler in separate goroutine
	go app.handleGracefulShutdown(srv, shutdownComplete)

	// Log server start with local IP for debugging
	localIP := app.getLocalIP()
	app.logger.Info("Starting HTTP server",
		"addr", srv.Addr,
		"local_ip", localIP,
		"shutdown_timeout", app.config.shutdownTimeout)

	// Start server - blocks until shutdown or error
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// Wait for graceful shutdown to complete
	<-shutdownComplete
	app.logger.Info("Server shutdown completed successfully")

	return nil
}

// handleGracefulShutdown listens for OS termination signals and performs
// controlled server shutdown within the configured timeout period.
func (app *application) handleGracefulShutdown(srv *http.Server, done chan<- struct{}) {
	// Create buffered channel for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	sig := <-quit
	app.logger.Info("Received shutdown signal", "signal", sig.String())

	// Create context with timeout for controlled shutdown
	ctx, cancel := context.WithTimeout(context.Background(), app.config.shutdownTimeout)
	defer cancel()

	// Attempt graceful server shutdown
	if err := srv.Shutdown(ctx); err != nil {
		app.logger.Error("Server forced to shutdown", "error", err)
	}

	// Signal shutdown completion
	close(done)
}
