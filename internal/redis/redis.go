package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisClient{
		Client: client,
	}, nil
}

func (r *RedisClient) LPush(ctx context.Context, key string, value string) error {
	err := r.Client.LPush(ctx, key, value).Err()
	if err != nil {
		return fmt.Errorf("failed to lpush: %w", err)
	}

	return nil
}

func (r *RedisClient) RPop(ctx context.Context, key string) (string, error) {
	value, err := r.Client.RPop(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", fmt.Errorf("failed to rpop: %w", err)
	}

	return value, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
