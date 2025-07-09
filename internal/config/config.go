package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	PostgresURL string `env:"POSTGRES_URL,required"`
	RedisURL    string `env:"REDIS_URL,required"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
