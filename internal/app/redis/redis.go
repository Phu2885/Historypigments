package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const servicePrefix = "historypigments."
const jwtPrefix = "jwt_blacklist."

type Client struct {
	rdb *redis.Client
}

func New(ctx context.Context, host, port, password string) (*Client, error) {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func getJWTKey(tokenStr string) string {
	return servicePrefix + jwtPrefix + tokenStr
}

// Запись отозванного JWT токена в Blacklist в Redis
func (c *Client) WriteJWTToBlacklist(ctx context.Context, tokenStr string, ttl time.Duration) error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Set(ctx, getJWTKey(tokenStr), "revoked", ttl).Err()
}

// Проверка, находится ли JWT токен в Blacklist в Redis
func (c *Client) CheckJWTInBlacklist(ctx context.Context, tokenStr string) bool {
	if c == nil || c.rdb == nil {
		return false
	}
	val, err := c.rdb.Get(ctx, getJWTKey(tokenStr)).Result()
	return err == nil && val == "revoked"
}

func (c *Client) Close() error {
	if c != nil && c.rdb != nil {
		return c.rdb.Close()
	}
	return nil
}
