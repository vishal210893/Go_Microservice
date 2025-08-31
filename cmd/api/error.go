// Package main provides HTTP error handling utilities for the API server.
package main

import (
	"log"
	"net/http"
)

// 500 — don't leak internals to clients; log full context server-side.
func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("500 %s %s: %v", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

// 400 — echo validation/parse message (already considered safe by caller).
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("400 %s %s: %v", r.Method, r.URL.Path, err)
	msg := "Bad Request"
	if err != nil && err.Error() != "" {
		msg = err.Error()
	}
	_ = writeJSONError(w, http.StatusBadRequest, msg)
}

// 404 — stable message; path is already in server logs above.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("404 %s %s: %v", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusNotFound, "Not Found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("409 %s %s: %v", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	err = writeJSONError(w, http.StatusUnauthorized, "unauthorized")
	if err != nil {
		return
	}
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	err = writeJSONError(w, http.StatusUnauthorized, "unauthorized")
	if err != nil {
		return
	}
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warn("403 %s %s", r.Method, r.URL.Path)
	_ = writeJSONError(w, http.StatusForbidden, "Forbidden")
}
