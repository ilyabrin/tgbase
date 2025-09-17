package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewRedisClient(ctx context.Context, cfg struct {
	Addr     string
	Password string
	DB       int
}) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration int64) error {
	if expiration > 0 {
		return c.client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	}
	// If expiration is 0, set without expiration
	return c.client.Set(ctx, key, value, 0).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	_, err := c.client.Del(ctx, key).Result()
	return err
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) HSet(ctx context.Context, key string, field string, value any) error {
	return c.client.HSet(ctx, key, field, value).Err()
}

func (c *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	val, err := c.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil // Field does not exist
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *Client) HDel(ctx context.Context, key string, field string) error {
	_, err := c.client.HDel(ctx, key, field).Result()
	return err
}

func (c *Client) HExists(ctx context.Context, key string, field string) (bool, error) {
	exists, err := c.client.HExists(ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	val, err := c.client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Key does not exist
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}
