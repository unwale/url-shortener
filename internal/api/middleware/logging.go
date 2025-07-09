package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	loggerKey contextKey = "logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqId := uuid.New().String()

		logger := slog.With(
			"request_id", reqId,
			"method", r.Method,
			"url", r.URL.String(),
		)
		ctx := context.WithValue(r.Context(), loggerKey, logger)

		logger.Info("Received request")
		next.ServeHTTP(w, r.WithContext(ctx))

		duration := time.Since(start)
		slog.Info("Processed request",
			"method", r.Method,
			"url", r.URL.String(),
			"duration", duration,
		)
	})
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
