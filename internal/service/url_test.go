package service

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	db "github.com/unwale/url-shortener/db/sqlc"
	"github.com/unwale/url-shortener/internal/domain/cache"
	"github.com/unwale/url-shortener/internal/domain/model"
	"github.com/unwale/url-shortener/internal/domain/repository"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateURL(ctx context.Context, params *db.CreateUrlParams) (*model.Url, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Url), args.Error(1)
}

func (m *mockRepository) GetURLByShortened(ctx context.Context, shortened string) (*model.Url, error) {
	args := m.Called(ctx, shortened)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Url), args.Error(1)
}

func (m *mockRepository) IncrementClickCount(ctx context.Context, shortened string) error {
	args := m.Called(ctx, shortened)
	return args.Error(0)
}

type mockCache struct {
	mock.Mock
}

func (m *mockCache) Get(ctx context.Context, key string) (*string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *mockCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func TestCreateShortURL_Success_WithAlias(t *testing.T) {
	mockRepo := new(mockRepository)
	mockCache := new(mockCache)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewURLService(mockRepo, mockCache, logger)

	originalURL := "https://www.google.com"
	alias := "my-google"

	expectedModel := &model.Url{OriginalUrl: originalURL, ShortUrl: alias}
	mockRepo.On("CreateURL", mock.Anything, &db.CreateUrlParams{
		OriginalUrl: originalURL,
		ShortUrl:    alias,
	}).Return(expectedModel, nil)

	shortURL, err := service.CreateShortURL(context.Background(), originalURL, alias)

	assert.NoError(t, err)
	assert.Equal(t, "my-google", shortURL)

	mockRepo.AssertExpectations(t)
}

func TestCreateURL_Success_NoAlias(t *testing.T) {
	mockRepo := new(mockRepository)
	mockCache := new(mockCache)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewURLService(mockRepo, mockCache, logger)

	originalURL := "https://www.google.com"

	expectedModel := &model.Url{OriginalUrl: originalURL, ShortUrl: "ac6bb669"}
	mockRepo.On("CreateURL", mock.Anything, &db.CreateUrlParams{
		OriginalUrl: originalURL,
		ShortUrl:    "ac6bb669",
	}).Return(expectedModel, nil)

	shortURL, err := service.CreateShortURL(context.Background(), originalURL, "")

	assert.NoError(t, err)
	assert.Equal(t, "ac6bb669", shortURL)

	mockRepo.AssertExpectations(t)
}

func TestCreateShortURL_Failure_AliasAlreadyExists(t *testing.T) {
	mockRepo := new(mockRepository)
	mockCache := new(mockCache)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewURLService(mockRepo, mockCache, logger)

	originalURL := "https://www.google.com"
	alias := "my-google"

	mockRepo.On("CreateURL", mock.Anything, &db.CreateUrlParams{
		OriginalUrl: originalURL,
		ShortUrl:    alias,
	}).Return(nil, repository.ErrURLAlreadyExists)

	shortURL, err := service.CreateShortURL(context.Background(), originalURL, alias)

	assert.Error(t, err)
	assert.Equal(t, "", shortURL)

	mockRepo.AssertExpectations(t)
}

func TestResolveShortURL_Success_CacheMiss(t *testing.T) {
	mockRepo := new(mockRepository)
	mockCache := new(mockCache)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewURLService(mockRepo, mockCache, logger)

	shortURL := "ac6bb669"
	originalURL := "https://www.google.com"

	mockCache.On("Get", mock.Anything, shortURL).Return(nil, cache.ErrCacheMiss)
	mockRepo.On("GetURLByShortened", mock.Anything, shortURL).Return(&model.Url{OriginalUrl: originalURL}, nil)
	mockRepo.On("IncrementClickCount", mock.Anything, shortURL).Return(nil)
	mockCache.On("Set", mock.Anything, shortURL, originalURL, CacheExpiration).Return(nil)

	resolvedURL, err := service.ResolveShortURL(context.Background(), shortURL)

	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, resolvedURL)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestResolveShortURL_Success_CacheHit(t *testing.T) {
	mockRepo := new(mockRepository)
	mockCache := new(mockCache)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewURLService(mockRepo, mockCache, logger)

	shortURL := "ac6bb669"
	originalURL := "https://www.google.com"

	mockCache.On("Get", mock.Anything, shortURL).Return(&originalURL, nil)
	mockRepo.On("IncrementClickCount", mock.Anything, shortURL).Return(nil)

	resolvedURL, err := service.ResolveShortURL(context.Background(), shortURL)

	time.Sleep(10 * time.Millisecond)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, resolvedURL)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)

	mockRepo.AssertNotCalled(t, "GetURLByShortened", mock.Anything, shortURL)
}

func TestGetShortURLStats_Success(t *testing.T) {
	mockRepo := new(mockRepository)
	mockCache := new(mockCache)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewURLService(mockRepo, mockCache, logger)

	shortURL := "ac6bb669"
	expectedStats := &model.Url{
		OriginalUrl: "https://www.google.com",
		ShortUrl:    shortURL,
		ClickCount:  10,
		CreatedAt:   "",
		UpdatedAt:   "",
	}

	mockRepo.On("GetURLByShortened", mock.Anything, shortURL).Return(expectedStats, nil)

	stats, err := service.GetShortURLStats(context.Background(), shortURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedStats, stats)

	mockRepo.AssertExpectations(t)
}

func TestGetShortURLStats_Failure_NotFound(t *testing.T) {
	mockRepo := new(mockRepository)
	mockCache := new(mockCache)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewURLService(mockRepo, mockCache, logger)

	shortURL := "ac6bb669"

	mockRepo.On("GetURLByShortened", mock.Anything, shortURL).Return(nil, repository.ErrURLNotFound)

	stats, err := service.GetShortURLStats(context.Background(), shortURL)

	assert.Error(t, err)
	assert.Nil(t, stats)

	mockRepo.AssertExpectations(t)
}
