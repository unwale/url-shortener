package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/unwale/url-shortener/db/sqlc"
	"github.com/unwale/url-shortener/internal/domain/model"
)

type URLRepository interface {
	CreateURL(ctx context.Context, url *db.CreateUrlParams) (*model.Url, error)
	GetURLByShortened(ctx context.Context, shortened string) (*model.Url, error)
	IncrementClickCount(ctx context.Context, shortened string) error
}

type urlRepository struct {
	querier db.Querier
}

func NewURLRepository(conn *pgxpool.Pool) URLRepository {
	return &urlRepository{
		querier: db.New(conn),
	}
}

func (r *urlRepository) CreateURL(ctx context.Context, url *db.CreateUrlParams) (*model.Url, error) {
	_, err := r.querier.GetUrlByShort(ctx, url.ShortUrl)
	if err == nil {
		return nil, ErrURLAlreadyExists
	}
	createdUrl, err := r.querier.CreateUrl(ctx,
		db.CreateUrlParams{
			OriginalUrl: url.OriginalUrl,
			ShortUrl:    url.ShortUrl,
		})

	return &model.Url{
		OriginalUrl: createdUrl.OriginalUrl,
		ShortUrl:    createdUrl.ShortUrl,
		CreatedAt:   createdUrl.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   createdUrl.UpdatedAt.Time.Format(time.RFC3339),
	}, err
}

func (r *urlRepository) GetURLByShortened(ctx context.Context, shortened string) (*model.Url, error) {
	url, err := r.querier.GetUrlByShort(ctx, shortened)
	if err != nil {
		return nil, ErrURLNotFound
	}

	return &model.Url{
		OriginalUrl: url.OriginalUrl,
		ShortUrl:    url.ShortUrl,
		ClickCount:  url.ClickCount,
		CreatedAt:   url.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   url.UpdatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (r *urlRepository) IncrementClickCount(ctx context.Context, shortened string) error {
	_, err := r.querier.IncrementClickCount(ctx, shortened)
	if err != nil {
		return ErrURLNotFound
	}
	return nil
}

var (
	ErrURLAlreadyExists = model.Error{
		Message: "URL already exists",
	}
	ErrURLNotFound = model.Error{
		Message: "URL not found",
	}
)
