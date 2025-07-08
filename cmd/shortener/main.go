package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/unwale/url-shortener/internal/api/handler"
	"github.com/unwale/url-shortener/internal/domain/repository"
	"github.com/unwale/url-shortener/internal/service"
)

func main() {
	context := context.Background()

	dbURL := os.Getenv("POSTGRES_URL")
	conn, err := pgxpool.New(context, dbURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	urlRepository := repository.NewURLRepository(conn)
	urlService := service.NewURLService(urlRepository)
	urlHandler := handler.NewURLHandler(urlService)

	mux := mux.NewRouter()
	urlHandler.RegisterRoutes(mux)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
