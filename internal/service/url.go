package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	db "github.com/unwale/url-shortener/db/sqlc"
	"github.com/unwale/url-shortener/internal/domain/model"
	"github.com/unwale/url-shortener/internal/domain/repository"
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL, alias string) (string, error)
	ResolveShortURL(ctx context.Context, shortURL string) (string, error)
}

type urlService struct {
	repository repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{
		repository: repo,
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
	url, err := s.repository.GetURLByShortened(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return url.OriginalUrl, nil
}

var (
	ErrInvalidAliasFormat = model.Error{
		Message: "Alias must be alphanumeric and between 3 to 20 characters long",
	}
	ErrAliasReserved = model.Error{
		Message: "Alias is reserved and cannot be used",
	}
)
