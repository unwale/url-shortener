package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/unwale/url-shortener/internal/api/handler"
	"github.com/unwale/url-shortener/internal/api/middleware"
	"github.com/unwale/url-shortener/internal/config"
	"github.com/unwale/url-shortener/internal/domain/cache"
	"github.com/unwale/url-shortener/internal/domain/repository"
	"github.com/unwale/url-shortener/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	slog.Info("Starting URL Shortener Service")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		panic(err)
	}

	ctx := context.Background()

	logger.Info("Connecting to PostgreSQL database")
	dbURL := cfg.PostgresURL
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer conn.Close()
	logger.Info("Connected to PostgreSQL database")

	logger.Info("Connecting to Redis cache")
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
		DB:   0, // use default DB
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("Failed to close Redis connection", "error", err)
		}
	}()
	logger.Info("Connected to Redis cache")

	urlCache := cache.NewRedisURLCache(redisClient)
	urlRepository := repository.NewURLRepository(conn)
	urlService := service.NewURLService(urlRepository, urlCache, *logger)
	urlHandler := handler.NewURLHandler(urlService)

	mux := mux.NewRouter()
	mux.Use(middleware.LoggingMiddleware)
	urlHandler.RegisterRoutes(mux)

	stopCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		logger.Info("Starting HTTP server on port 8080")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start HTTP server", "error", err)
			os.Exit(1)
		}
	}()

	<-stopCtx.Done()

	logger.Info("Shutting down server gracefully")
	timeoutCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(timeoutCtx); err != nil {
		logger.Error("Failed to shutdown HTTP server", "error", err)
	} else {
		logger.Info("HTTP server shut down gracefully")
	}
}
