package main

import (
	"Go-Microservice/docs"
	"Go-Microservice/internal/db"
	"Go-Microservice/internal/env"
	formatLog "Go-Microservice/internal/log"
	"Go-Microservice/internal/repo"
	"context"
	"log"
	"log/slog"
	"net"
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

// @title						Go Microservice API
// @version					1.0
// @description				A production-ready Go microservice with posts, users, and social features
// @description				This API provides endpoints for managing posts, users, comments, and social interactions
// @description				including user following/unfollowing and personalized feeds.
//
// @contact.name				API Support
// @contact.email				support@example.com
//
// @license.name				MIT
// @license.url				https://opensource.org/licenses/MIT
//
// @host						localhost:8080
// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
//
// @tag.name				posts
// @tag.description		Operations related to posts management
//
// @tag.name				users
// @tag.description		Operations related to user management and social features
//
// @tag.name				health
// @tag.description		Health check and system status endpoints
//
// @schemes				http https
// @produce				json
// @consumes				json
//
// @x-extension-openapi	{"info":{"x-logo":{"url":"https://example.com/logo.png"}}}
func main() {

	logger := slog.New(formatLog.NewFormattedLogHandler(os.Stdout, slog.LevelInfo))
	slog.SetDefault(logger)

	config := config{
		addr:            env.GetPort("ADDR", DefaultAddr),
		shutdownTimeout: DefaultServerTimeout,
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://avnadmin:AVNS_LT5DsEKUPKfrHSHZHyB@pg-1d9d15dc-vishal210893-5985.h.aivencloud.com:28832/defaultdb?sslmode=require"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		apiUrl: env.GetString("API_URL", "localhost:8000"),
	}

	db, err := db.New(
		config.db.addr,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Println("Connected to database!")

	postgresRepo, _ := repo.NewPostgresRepo(db)

	app := &application{
		config: config,
		logger: logger,
		repo:   *postgresRepo,
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
	//Docs
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = app.config.apiUrl
	docs.SwaggerInfo.BasePath = ""

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
