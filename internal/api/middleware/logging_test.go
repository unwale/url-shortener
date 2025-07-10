package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler := LoggingMiddleware(mockHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	logger := GetLoggerFromContext(req.Context())
	assert.NotNil(t, logger)
}

func TestGetLoggerFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), loggerKey, &slog.Logger{})
	logger := GetLoggerFromContext(ctx)

	assert.NotNil(t, logger)
	assert.IsType(t, &slog.Logger{}, logger)

	ctx = context.Background()
	logger = GetLoggerFromContext(ctx)
	assert.Equal(t, slog.Default(), logger)
}
