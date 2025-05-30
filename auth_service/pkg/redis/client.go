package redis

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("key not found")
)

type Client struct {
	cli *redis.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cli: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
	}
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.cli.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.cli.Get(ctx, key).Result()
}

func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.cli.Del(ctx, keys...).Err()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.cli.Ping(ctx).Err()
}
