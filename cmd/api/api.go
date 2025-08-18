package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net"
	"net/http"
	"time"
)

type config struct {
	addr string
}

type application struct {
	config config
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthcheckHandler)
	})

	r.Get("/v1/health", app.healthcheckHandler)
	return r
}

func (app *application) run(mux http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	ip := getLocalIP()
	log.Printf("Starting server at %s (IP: %s)", app.config.addr, ip)
	return server.ListenAndServe()
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("failed to get local IP addresses: %v", err)
		return "Unknown"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "Unknown"
}
