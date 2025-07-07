package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	db "github.com/unwale/url-shortener/db/sqlc"
	"github.com/unwale/url-shortener/internal/domain/model"
)

type URLRepository interface {
	CreateURL(ctx context.Context, url *db.CreateUrlParams) (*model.Url, error)
	GetURLByShortened(ctx context.Context, shortened string) (*model.Url, error)
}

type urlRepository struct {
	querier db.Querier
}

func NewURLRepository(conn *pgx.Conn) URLRepository {
	return &urlRepository{
		querier: db.New(conn),
	}
}

func (r *urlRepository) CreateURL(ctx context.Context, url *db.CreateUrlParams) (*model.Url, error) {
	_, err := r.querier.GetUrlByShort(ctx, url.ShortUrl)
	if err != nil {
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
		CreatedAt:   url.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   url.UpdatedAt.Time.Format(time.RFC3339),
	}, nil
}

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrURLAlreadyExists = Error{
		Message: "URL already exists",
	}
	ErrURLNotFound = Error{
		Message: "URL not found",
	}
)
