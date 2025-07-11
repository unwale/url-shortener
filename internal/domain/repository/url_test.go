//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/unwale/url-shortener/db/sqlc"
	"github.com/unwale/url-shortener/internal/config"
)

var (
	testPool *pgxpool.Pool
)

func runWithTestDb(t *testing.T, fn func(repo *URLRepository)) {
	t.Cleanup(func() {
		_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE urls")
		require.NoError(t, err)
	})

	if testPool == nil {
		cfg, err := config.LoadConfig()
		require.NoError(t, err)
		testPool, err = pgxpool.New(context.Background(), cfg.PostgresURL)
		require.NoError(t, err)
	}

	repo := NewURLRepository(testPool)
	fn(&repo)
}

func TestCreateURL(t *testing.T) {
	t.Run("create url", func(t *testing.T) {
		runWithTestDb(t, func(repo *URLRepository) {
			url := &db.CreateUrlParams{
				OriginalUrl: "https://google.com",
				ShortUrl:    "exmpl",
			}

			createdURL, err := (*repo).CreateURL(context.Background(), url)

			require.NoError(t, err)
			assert.NotNil(t, createdURL)
			assert.Equal(t, url.OriginalUrl, createdURL.OriginalUrl)
			assert.Equal(t, url.ShortUrl, createdURL.ShortUrl)
			assert.NotEmpty(t, createdURL.CreatedAt)
			assert.NotEmpty(t, createdURL.UpdatedAt)
		})
	})

	t.Run("create duplicate url", func(t *testing.T) {
		runWithTestDb(t, func(repo *URLRepository) {
			url := &db.CreateUrlParams{
				OriginalUrl: "https://google.com",
				ShortUrl:    "exmpl",
			}

			_, err := (*repo).CreateURL(context.Background(), url)
			require.NoError(t, err)

			_, err = (*repo).CreateURL(context.Background(), url)
			assert.ErrorIs(t, err, ErrURLAlreadyExists)
		})
	})
}

func TestGetURLByShortened(t *testing.T) {
	t.Run("get existing url", func(t *testing.T) {
		runWithTestDb(t, func(repo *URLRepository) {
			url := &db.CreateUrlParams{
				OriginalUrl: "https://google.com",
				ShortUrl:    "exmpl",
			}

			createdURL, err := (*repo).CreateURL(context.Background(), url)
			require.NoError(t, err)

			fetchedURL, err := (*repo).GetURLByShortened(context.Background(), createdURL.ShortUrl)
			require.NoError(t, err)
			assert.NotNil(t, fetchedURL)
			assert.Equal(t, createdURL.OriginalUrl, fetchedURL.OriginalUrl)
			assert.Equal(t, createdURL.ShortUrl, fetchedURL.ShortUrl)
		})
	})

	t.Run("get non-existing url", func(t *testing.T) {
		runWithTestDb(t, func(repo *URLRepository) {
			fetchedURL, err := (*repo).GetURLByShortened(context.Background(), "nonexistent")
			assert.ErrorIs(t, err, ErrURLNotFound)
			assert.Nil(t, fetchedURL)
		})
	})
}

func TestIncrementClickCount(t *testing.T) {
	t.Run("increment click count", func(t *testing.T) {
		runWithTestDb(t, func(repo *URLRepository) {
			url := &db.CreateUrlParams{
				OriginalUrl: "https://google.com",
				ShortUrl:    "exmpl",
			}

			createdURL, err := (*repo).CreateURL(context.Background(), url)
			require.NoError(t, err)

			err = (*repo).IncrementClickCount(context.Background(), createdURL.ShortUrl)
			require.NoError(t, err)

			fetchedURL, err := (*repo).GetURLByShortened(context.Background(), createdURL.ShortUrl)
			require.NoError(t, err)
			assert.Equal(t, 1, int(fetchedURL.ClickCount))
		})
	})

	t.Run("increment click count for non-existing url", func(t *testing.T) {
		runWithTestDb(t, func(repo *URLRepository) {
			err := (*repo).IncrementClickCount(context.Background(), "nonexistent")
			assert.ErrorIs(t, err, ErrURLNotFound)
		})
	})
}
