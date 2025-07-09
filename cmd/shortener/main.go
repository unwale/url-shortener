package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/unwale/url-shortener/internal/api/handler"
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

	context := context.Background()

	logger.Info("Connecting to PostgreSQL database")
	dbURL := os.Getenv("POSTGRES_URL")
	conn, err := pgxpool.New(context, dbURL)
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", "error", err)
		panic(err)
	}
	defer conn.Close()
	logger.Info("Connected to PostgreSQL database")

	logger.Info("Connecting to Redis cache")
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
		DB:   0, // use default DB
	})
	if err := redisClient.Ping(context).Err(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}
	defer redisClient.Close()
	logger.Info("Connected to Redis cache")

	urlCache := cache.NewRedisURLCache(redisClient)
	urlRepository := repository.NewURLRepository(conn)
	urlService := service.NewURLService(urlRepository, urlCache)
	urlHandler := handler.NewURLHandler(urlService)

	mux := mux.NewRouter()
	urlHandler.RegisterRoutes(mux)

	logger.Info("Starting HTTP server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("Failed to start HTTP server", "error", err)
		panic(err)
	}
}
