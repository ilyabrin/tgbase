package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RealRedisClient wraps the actual Redis client to implement RedisInterface
type RealRedisClient struct {
	client *redis.Client
}

// NewRealRedisClient creates a new real Redis client wrapper
func NewRealRedisClient(ctx context.Context, cfg Config) (*RealRedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RealRedisClient{client: client}, nil
}

func (r *RealRedisClient) Set(ctx context.Context, key string, value any, expiration int64) error {
	if expiration > 0 {
		return r.client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	}
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RealRedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RealRedisClient) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func (r *RealRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RealRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *RealRedisClient) Decr(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *RealRedisClient) HSet(ctx context.Context, key string, field string, value any) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

func (r *RealRedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	val, err := r.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RealRedisClient) HDel(ctx context.Context, key string, field string) error {
	_, err := r.client.HDel(ctx, key, field).Result()
	return err
}

func (r *RealRedisClient) HExists(ctx context.Context, key string, field string) (bool, error) {
	exists, err := r.client.HExists(ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *RealRedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	val, err := r.client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return make(map[string]string), nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *RealRedisClient) Close() error {
	return r.client.Close()
}
