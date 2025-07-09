package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"strings"
	"time"

	db "github.com/unwale/url-shortener/db/sqlc"
	"github.com/unwale/url-shortener/internal/domain/cache"
	"github.com/unwale/url-shortener/internal/domain/model"
	"github.com/unwale/url-shortener/internal/domain/repository"
)

const (
	CacheExpiration = 24 * time.Hour
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL, alias string) (string, error)
	ResolveShortURL(ctx context.Context, shortURL string) (string, error)
	GetShortURLStats(ctx context.Context, shortURL string) (*model.Url, error)
}

type urlService struct {
	repository repository.URLRepository
	cache      cache.URLCache
	logger     *slog.Logger
}

func NewURLService(repo repository.URLRepository, cache cache.URLCache, logger *slog.Logger) URLService {
	return &urlService{
		repository: repo,
		cache:      cache,
		logger:     logger,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, originalURL, alias string) (string, error) {
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	var shortURL string

	if alias != "" {
		if len(alias) < 4 || len(alias) > 20 {
			return "", ErrInvalidAliasFormat
		}
		if strings.HasPrefix(alias, "api/") {
			return "", ErrAliasReserved
		}
		shortURL = alias
	} else {

		hash := sha256.Sum256([]byte(originalURL))
		shortURL = hex.EncodeToString(hash[:])[:8]
	}

	model, err := s.repository.CreateURL(ctx, &db.CreateUrlParams{
		OriginalUrl: originalURL,
		ShortUrl:    shortURL})
	if err != nil {
		return "", err
	}
	return model.ShortUrl, nil
}

func (s *urlService) ResolveShortURL(ctx context.Context, shortURL string) (string, error) {
	originalUrl, err := s.cache.Get(ctx, shortURL)
	if err == nil {
		go func() {
			backgroundCtx := context.Background()
			if err := s.repository.IncrementClickCount(backgroundCtx, shortURL); err != nil {
				s.logger.Error("Failed to increment click count", "shortURL", shortURL, "error", err)
			}
		}()
		return *originalUrl, nil
	}

	url, err := s.repository.GetURLByShortened(ctx, shortURL)
	if err != nil {
		return "", err
	}

	go func() {
		backgroundCtx := context.Background()
		if err := s.repository.IncrementClickCount(backgroundCtx, shortURL); err != nil {
			s.logger.Error("Failed to increment click count", "shortURL", shortURL, "error", err)
		}
		if err := s.cache.Set(backgroundCtx, shortURL, url.OriginalUrl, CacheExpiration); err != nil {
			s.logger.Error("Failed to cache URL", "shortURL", shortURL, "error", err)
		}
	}()

	return url.OriginalUrl, nil
}

func (s *urlService) GetShortURLStats(ctx context.Context, shortURL string) (*model.Url, error) {
	url, err := s.repository.GetURLByShortened(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	return &model.Url{
		OriginalUrl: url.OriginalUrl,
		ShortUrl:    url.ShortUrl,
		ClickCount:  url.ClickCount,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
	}, nil
}

var (
	ErrInvalidAliasFormat = model.Error{
		Message: "Alias must be alphanumeric and between 4 to 20 characters long",
	}
	ErrAliasReserved = model.Error{
		Message: "Alias is reserved and cannot be used",
	}
)
