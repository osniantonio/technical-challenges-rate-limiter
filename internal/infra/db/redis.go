package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
	"github.com/redis/go-redis/v9"
)

type RedisDatabaseGateway struct {
	client *redis.Client
}

func newRedisDatabaseGateway(options *gateway.DatabaseOptions) (gateway.DatabaseGateway, error) {
	db, _ := strconv.Atoi(options.Database)
	addr := fmt.Sprintf("%s:%s", options.Host, options.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: options.Password,
		DB:       db,
	})
	return &RedisDatabaseGateway{client}, nil
}

func (g *RedisDatabaseGateway) Lock(ctx context.Context, key string, ttl time.Duration) error {
	key = fmt.Sprintf("lock:%s", key)
	if _, err := g.client.SetNX(ctx, key, true, ttl).Result(); err != nil {
		return err
	}
	return nil
}

func (g *RedisDatabaseGateway) IsLocked(ctx context.Context, key string) (bool, error) {
	key = fmt.Sprintf("lock:%s", key)
	_, err := g.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *RedisDatabaseGateway) SaveRequest(ctx context.Context, key string) error {
	key = fmt.Sprintf("request:%s", key)
	if _, err := g.client.SetNX(ctx, key, 0, time.Second).Result(); err != nil {
		return err
	}
	if _, err := g.client.Incr(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}

func (g *RedisDatabaseGateway) CountRequests(ctx context.Context, key string) (int, error) {
	key = fmt.Sprintf("request:%s", key)
	result, err := g.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	value, _ := strconv.Atoi(result)
	return value, nil
}

func (g *RedisDatabaseGateway) CreateToken(ctx context.Context, token string, limit int) error {
	key := fmt.Sprintf("limit:%s", token)
	if _, err := g.client.Set(ctx, key, limit, 0).Result(); err != nil {
		return err
	}
	return nil
}

func (g *RedisDatabaseGateway) GetTokenLimit(ctx context.Context, token string) (int, error) {
	key := fmt.Sprintf("limit:%s", token)
	result, err := g.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return 0, nil
	}
	if result == "" {
		return 0, nil
	}
	limit, _ := strconv.Atoi(result)
	return limit, nil
}
