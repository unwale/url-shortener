package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLCache interface {
	Get(ctx context.Context, key string) (*string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type RedisURLCache struct {
	client *redis.Client
}

func NewRedisURLCache(client *redis.Client) URLCache {
	return &RedisURLCache{
		client: client,
	}
}

func (c *RedisURLCache) Get(ctx context.Context, key string) (*string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &val, nil
}

func (c *RedisURLCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := c.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
