package redis

import (
	"context"
)

// Config holds Redis connection parameters.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// RedisInterface defines the contract for Redis operations
type RedisInterface interface {
	Set(ctx context.Context, key string, value any, expiration int64) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	HSet(ctx context.Context, key string, field string, value any) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HDel(ctx context.Context, key string, field string) error
	HExists(ctx context.Context, key string, field string) (bool, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Close() error
}

type Client struct {
	client RedisInterface
}

// NewRedisClient creates a new Redis client with a real Redis connection.
func NewRedisClient(ctx context.Context, cfg Config) (*Client, error) {
	realClient, err := NewRealRedisClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Client{client: realClient}, nil
}

// NewMockClient creates a Redis client backed by an in-memory mock (for testing).
func NewMockClient() *Client {
	return &Client{client: NewMockRedisClient()}
}

func (c *Client) Close() error { return c.client.Close() }

func (c *Client) Set(ctx context.Context, key string, value any, expiration int64) error {
	return c.client.Set(ctx, key, value, expiration)
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key)
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key)
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	return c.client.Exists(ctx, key)
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key)
}

func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key)
}

func (c *Client) HSet(ctx context.Context, key string, field string, value any) error {
	return c.client.HSet(ctx, key, field, value)
}

func (c *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	return c.client.HGet(ctx, key, field)
}

func (c *Client) HDel(ctx context.Context, key string, field string) error {
	return c.client.HDel(ctx, key, field)
}

func (c *Client) HExists(ctx context.Context, key string, field string) (bool, error) {
	return c.client.HExists(ctx, key, field)
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key)
}
