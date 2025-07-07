package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/unwale/url-shortener/internal/domain/model"
	"github.com/unwale/url-shortener/internal/domain/repository"
)

type URLService interface {
	GenerateShortURL(originalURL string) (string, error)
	ResolveShortURL(shortURL string) (string, error)
}

type urlService struct {
	repository repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{
		repository: repo,
	}
}

func (s *urlService) GenerateShortURL(originalURL string) (string, error) {
	hash := sha256.Sum256([]byte(originalURL))
	shortURL := hex.EncodeToString(hash[:])[:8]

	model, err := s.repository.CreateURL(&model.URL{
		ID:        shortURL,
		Original:  originalURL,
		Shortened: shortURL,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})

	return model.Shortened, err
}

func (s *urlService) ResolveShortURL(shortURL string) (string, error) {
	url, err := s.repository.GetURLByShortened(shortURL)
	if err != nil {
		return "", err
	}
	return url.Original, nil
}
