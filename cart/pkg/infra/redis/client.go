package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"route256/cart/internal/config"
)

type Client struct {
	client *redis.Client
	TTL    time.Duration
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Address.Host, cfg.Redis.Address.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
		TTL: cfg.Redis.TTL,
	}
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Client) Close() error {
	return c.client.Close()
}
