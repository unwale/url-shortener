package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	db "github.com/unwale/url-shortener/db/sqlc"
	"github.com/unwale/url-shortener/internal/domain/repository"
)

type URLService interface {
	GenerateShortURL(ctx context.Context, originalURL string) (string, error)
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

func (s *urlService) GenerateShortURL(ctx context.Context, originalURL string) (string, error) {
	hash := sha256.Sum256([]byte(originalURL))
	shortURL := hex.EncodeToString(hash[:])[:8]

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
