package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/shownest/pkg/config"
)

// Config holds the Redis connection parameters.
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Init loads the Redis config from the provider and returns a Redis client.
func Init(ctx context.Context, provider config.ConfigProvider) (*redis.Client, error) {
	raw, err := provider.Get(ctx, config.RedisCredentials)
	if err != nil {
		return nil, fmt.Errorf("cache: get config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("cache: parse config: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("cache: ping: %w", err)
	}

	return client, nil
}
