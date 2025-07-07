package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"

	"github.com/unwale/url-shortener/internal/api/handler"
	"github.com/unwale/url-shortener/internal/domain/repository"
	"github.com/unwale/url-shortener/internal/service"
)

func main() {
	context := context.Background()
	conn, err := pgx.Connect(context, "postgres://user:password@localhost:5432/url_shortener")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context)

	urlRepository := repository.NewURLRepository(conn)
	urlService := service.NewURLService(urlRepository)
	urlHandler := handler.NewURLHandler(urlService)

	mux := mux.NewRouter()
	urlHandler.RegisterRoutes(mux)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
