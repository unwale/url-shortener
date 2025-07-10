package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/unwale/url-shortener/internal/api/handler"
	"github.com/unwale/url-shortener/internal/api/model"
	domain "github.com/unwale/url-shortener/internal/domain/model"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) CreateShortURL(ctx context.Context, originalURL string, alias string) (string, error) {
	args := m.Called(ctx, originalURL, alias)
	return args.String(0), args.Error(1)
}

func (m *MockURLService) ResolveShortURL(ctx context.Context, shortenedURL string) (string, error) {
	args := m.Called(ctx, shortenedURL)
	return args.String(0), args.Error(1)
}

func (m *MockURLService) GetShortURLStats(ctx context.Context, shortenedURL string) (*domain.Url, error) { // Assuming model is from domain
	args := m.Called(ctx, shortenedURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Url), args.Error(1)
}

func TestShortenURLHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		requestBody := `{"url":"https://google.com","alias":"my-google"}`
		expectedShortURL := "123xyz"

		mockService.On("CreateShortURL", mock.Anything, "https://google.com", "my-google").Return(expectedShortURL, nil)

		req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(requestBody))
		rr := httptest.NewRecorder()

		urlHandler.ShortenURLHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response model.ShortenURLResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedShortURL, response.ShortURL)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		requestBody := `invalid json`

		req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(requestBody))
		rr := httptest.NewRecorder()

		urlHandler.ShortenURLHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertNotCalled(t, "CreateShortURL")
	})

	t.Run("service returns error", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)
		requestBody := `{"url":"https://google.com"}`

		mockService.On("CreateShortURL", mock.Anything, "https://google.com", "").Return("", errors.New("something went wrong"))

		req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(requestBody))
		rr := httptest.NewRecorder()

		urlHandler.ShortenURLHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

}

func TestResolveShortURLHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		shortened := "123xyz"
		originalURL := "https://google.com"

		mockService.On("ResolveShortURL", mock.Anything, shortened).Return(originalURL, nil)

		req := httptest.NewRequest("GET", "/"+shortened, nil)
		req = mux.SetURLVars(req, map[string]string{"shortened": shortened})

		rr := httptest.NewRecorder()

		urlHandler.ResolveShortURLHandler(rr, req)

		assert.Equal(t, http.StatusPermanentRedirect, rr.Code)
		assert.Equal(t, originalURL, rr.Header().Get("Location"))
		mockService.AssertExpectations(t)
	})

	t.Run("shortened URL empty", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		urlHandler.ResolveShortURLHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertNotCalled(t, "ResolveShortURL")
	})

	t.Run("service returns error", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		shortened := "123xyz"

		mockService.On("ResolveShortURL", mock.Anything, shortened).Return("", errors.New("not found"))

		req := httptest.NewRequest("GET", "/"+shortened, nil)
		req = mux.SetURLVars(req, map[string]string{"shortened": shortened})

		rr := httptest.NewRecorder()

		urlHandler.ResolveShortURLHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

}

func TestStatsHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		shortened := "123xyz"
		stats := &domain.Url{
			ShortUrl:    shortened,
			OriginalUrl: "https://google.com",
			ClickCount:  10,
		}

		mockService.On("GetShortURLStats", mock.Anything, shortened).Return(stats, nil)

		req := httptest.NewRequest("GET", "/api/stats/"+shortened, nil)
		req = mux.SetURLVars(req, map[string]string{"shortened": shortened})

		rr := httptest.NewRecorder()

		urlHandler.StatsHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response model.ShortUrlStatsResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 10, response.ClickCount)
		assert.Equal(t, "https://google.com", response.OriginalURL)
		mockService.AssertExpectations(t)
	})

	t.Run("shortened URL empty", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		req := httptest.NewRequest("GET", "/api/stats/", nil)
		rr := httptest.NewRecorder()

		urlHandler.StatsHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertNotCalled(t, "GetShortURLStats")
	})

	t.Run("service returns error", func(t *testing.T) {
		mockService := new(MockURLService)
		urlHandler := handler.NewURLHandler(mockService)

		shortened := "123xyz"

		mockService.On("GetShortURLStats", mock.Anything, shortened).Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/api/stats/"+shortened, nil)
		req = mux.SetURLVars(req, map[string]string{"shortened": shortened})

		rr := httptest.NewRecorder()

		urlHandler.StatsHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}
