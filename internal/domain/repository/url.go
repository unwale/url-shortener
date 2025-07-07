package repository

import (
	"github.com/unwale/url-shortener/internal/domain/model"
)

type URLRepository interface {
	CreateURL(url *model.URL) (*model.URL, error)
	GetURLByShortened(shortened string) (*model.URL, error)
	GetURLByID(id string) (*model.URL, error)
	UpdateURL(url *model.URL) (*model.URL, error)
	DeleteURL(id string) error
}

type urlRepository struct {
	urls map[string]*model.URL
}

func NewURLRepository() URLRepository {
	return &urlRepository{
		urls: make(map[string]*model.URL),
	}
}

func (r *urlRepository) CreateURL(url *model.URL) (*model.URL, error) {
	if _, exists := r.urls[url.Shortened]; exists {
		return nil, ErrURLAlreadyExists
	}
	r.urls[url.Shortened] = url
	return url, nil
}

func (r *urlRepository) GetURLByShortened(shortened string) (*model.URL, error) {
	url, exists := r.urls[shortened]
	if !exists {
		return nil, ErrURLNotFound
	}
	return url, nil
}

func (r *urlRepository) GetURLByID(id string) (*model.URL, error) {
	for _, url := range r.urls {
		if url.ID == id {
			return url, nil
		}
	}
	return nil, ErrURLNotFound
}

func (r *urlRepository) UpdateURL(url *model.URL) (*model.URL, error) {
	if _, exists := r.urls[url.Shortened]; !exists {
		return nil, ErrURLNotFound
	}
	r.urls[url.Shortened] = url
	return url, nil
}

func (r *urlRepository) DeleteURL(id string) error {
	for shortened, url := range r.urls {
		if url.ID == id {
			delete(r.urls, shortened)
			return nil
		}
	}
	return ErrURLNotFound
}

var (
	ErrURLAlreadyExists = model.Error{
		Message: "URL already exists",
	}
	ErrURLNotFound = model.Error{
		Message: "URL not found",
	}
)
