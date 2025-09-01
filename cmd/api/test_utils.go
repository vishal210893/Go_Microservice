package main

import (
	"Go-Microservice/internal/auth"
	formatLog "Go-Microservice/internal/log"
	"Go-Microservice/internal/repo"
	"Go-Microservice/internal/repo/cache"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	mockStore := repo.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	logger := slog.New(formatLog.NewFormattedLogHandler(os.Stdout, slog.LevelInfo))
	slog.SetDefault(logger)

	return &application{
		repo:          mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
		logger:        logger,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
