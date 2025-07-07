package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/unwale/url-shortener/internal/api/model"
	"github.com/unwale/url-shortener/internal/service"
)

type URLHandler struct {
	service service.URLService
}

func NewURLHandler(s service.URLService) *URLHandler {
	return &URLHandler{
		service: s,
	}
}

func (h *URLHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/shorten", h.ShortenURLHandler).Methods("POST")
	router.HandleFunc("/{shortened}", h.ResolveShortURLHandler).Methods("GET")
}

func (h *URLHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var request model.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shornetedURL, err := h.service.GenerateShortURL(request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.ShortenURLResponse{
		ShortURL: shornetedURL,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) ResolveShortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortened := vars["shortened"]
	if len(shortened) == 0 {
		http.Error(w, "Shortened URL is required", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.ResolveShortURL(shortened)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusPermanentRedirect)
}
