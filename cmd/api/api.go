package main

import (
	"Go-Microservice/internal/auth"
	"Go-Microservice/internal/mailer"
	"Go-Microservice/internal/repo"
	"Go-Microservice/internal/repo/cache"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// config holds all configuration parameters for the application
type config struct {
	// addr is the network address the server will listen on
	addr string
	// shutdownTimeout defines the maximum time allowed for graceful shutdown
	shutdownTimeout   time.Duration
	db                dbConfig
	apiUrl            string
	invitationExpTime time.Duration
	mailConfig        mailConfig
	frontendURL       string
	env               string
	auth              authConfig
	redisConfig       redisConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	aud    string
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
}

type sendGridConfig struct {
	apiKey string
}

// application holds the dependencies for HTTP handlers, helpers, and middleware.
// It serves as the dependency injection container for the entire application.
type application struct {
	config        config
	logger        *slog.Logger
	repo          repo.Repository
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
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
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthcheckHandler)

		// Swagger documentation endpoint
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Post("/comments", app.createCommentHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				//r.Use(app.usersContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}
