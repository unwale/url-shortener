package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		os.Setenv("POSTGRES_URL", "postgres://user:password@localhost:5432/dbname") //nolint:errcheck
		os.Setenv("REDIS_URL", "redis://localhost:6379/0")                          //nolint:errcheck

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.PostgresURL)
		assert.Equal(t, "redis://localhost:6379/0", cfg.RedisURL)
	})

	t.Run("missing required env vars", func(t *testing.T) {
		os.Unsetenv("POSTGRES_URL") //nolint:errcheck
		os.Unsetenv("REDIS_URL")    //nolint:errcheck

		cfg, err := LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)

		assert.Contains(t, err.Error(), "required environment variable")
	})
}
