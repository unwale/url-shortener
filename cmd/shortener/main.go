package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/unwale/url-shortener/internal/api/handler"
	"github.com/unwale/url-shortener/internal/domain/repository"
	"github.com/unwale/url-shortener/internal/service"
)

func main() {
	urlRepository := repository.NewURLRepository()
	urlService := service.NewURLService(urlRepository)
	urlHandler := handler.NewURLHandler(urlService)

	mux := mux.NewRouter()
	urlHandler.RegisterRoutes(mux)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
