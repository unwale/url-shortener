//go:build integration

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/redis/go-redis/v9"
	"github.com/unwale/url-shortener/internal/config"
)

var (
	testRedisClient *redis.Client
)

func runWithTestCache(t *testing.T, fn func(repo *URLCache)) {
	t.Cleanup(func() {
		err := testRedisClient.FlushDB(context.Background()).Err()
		require.NoError(t, err)
	})

	if testRedisClient == nil {
		cfg, err := config.LoadConfig()
		require.NoError(t, err)
		testRedisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisURL,
			Password: "",
			DB:       0,
		})
	}

	repo := NewRedisURLCache(testRedisClient)
	fn(&repo)
}

func TestSet(t *testing.T) {
	expiration := 5 * time.Second
	t.Run("set url", func(t *testing.T) {
		runWithTestCache(t, func(repo *URLCache) {
			err := (*repo).Set(context.Background(), "exmpl", "https://google.com", expiration)
			require.NoError(t, err)
		})
	})

	t.Run("set duplicate url", func(t *testing.T) {
		runWithTestCache(t, func(repo *URLCache) {
			err := (*repo).Set(context.Background(), "exmpl", "https://google.com", expiration)
			require.NoError(t, err)

			err = (*repo).Set(context.Background(), "exmpl", "https://google.com", expiration)
			require.NoError(t, err)
		})
	})
}

func TestGet(t *testing.T) {
	t.Run("get existing url", func(t *testing.T) {
		runWithTestCache(t, func(repo *URLCache) {
			err := (*repo).Set(context.Background(), "exmpl", "https://google.com", 5*time.Second)
			require.NoError(t, err)

			val, err := (*repo).Get(context.Background(), "exmpl")
			require.NoError(t, err)
			require.NotNil(t, val)
			require.Equal(t, "https://google.com", *val)
		})
	})

	t.Run("get non-existing url", func(t *testing.T) {
		runWithTestCache(t, func(repo *URLCache) {
			val, err := (*repo).Get(context.Background(), "non-existing")
			require.ErrorIs(t, err, ErrCacheMiss)
			require.Nil(t, val)
		})
	})
}
