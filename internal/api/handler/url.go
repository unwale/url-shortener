package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/unwale/url-shortener/internal/api/middleware"
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
	router.HandleFunc("/api/shorten", h.ShortenURLHandler).Methods("POST")
	router.HandleFunc("/{shortened}", h.ResolveShortURLHandler).Methods("GET")
	router.HandleFunc("/api/stats/{shortened}", h.StatsHandler).Methods("GET")
}

func (h *URLHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLoggerFromContext(r.Context())

	var request model.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shornetedURL, err := h.service.CreateShortURL(r.Context(), request.URL, request.Alias)
	if err != nil {
		logger.Error("Failed to create short URL", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.ShortenURLResponse{
		ShortURL: shornetedURL,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) ResolveShortURLHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLoggerFromContext(r.Context())

	vars := mux.Vars(r)
	shortened := vars["shortened"]
	if len(shortened) == 0 {
		logger.Error("Shortened URL is required")
		http.Error(w, "Shortened URL is required", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.ResolveShortURL(r.Context(), shortened)
	if err != nil {
		logger.Error("Failed to resolve short URL", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Redirecting to original URL", "shortened", shortened, "originalURL", originalURL)
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusPermanentRedirect)
}

func (h *URLHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLoggerFromContext(r.Context())

	vars := mux.Vars(r)
	shortened := vars["shortened"]
	if len(shortened) == 0 {
		logger.Error("Shortened URL is required")
		http.Error(w, "Shortened URL is required", http.StatusBadRequest)
		return
	}

	stats, err := h.service.GetShortURLStats(r.Context(), shortened)
	if err != nil {
		logger.Error("Failed to get short URL stats", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.ShortUrlStatsResponse{
		ShortURL:    stats.ShortUrl,
		OriginalURL: stats.OriginalUrl,
		ClickCount:  int(stats.ClickCount),
		CreatedAt:   stats.CreatedAt,
		UpdatedAt:   stats.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
