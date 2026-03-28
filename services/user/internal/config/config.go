package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pkgconfig "github.com/shownest/pkg/config"
)

type Config struct {
	Port             string
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
}

type rawConfig struct {
	Port             string `json:"port"`
	JWTAccessSecret  string `json:"jwtAccessSecret"`
	JWTRefreshSecret string `json:"jwtRefreshSecret"`
	JWTAccessExpiry  string `json:"jwtAccessExpiry"`
	JWTRefreshExpiry string `json:"jwtRefreshExpiry"`
}

// Load retrieves and parses the service config from the provider.
func Load(ctx context.Context, provider pkgconfig.ConfigProvider) (*Config, error) {
	raw, err := provider.Get(ctx, "service")
	if err != nil {
		return nil, fmt.Errorf("config: get service config: %w", err)
	}

	var rc rawConfig
	if err := json.Unmarshal(raw, &rc); err != nil {
		return nil, fmt.Errorf("config: parse service config: %w", err)
	}

	accessExpiry, err := time.ParseDuration(rc.JWTAccessExpiry) // access token expire: 15m
	if err != nil {
		return nil, fmt.Errorf("config: parse jwt_access_expiry %q: %w", rc.JWTAccessExpiry, err)
	}
	refreshExpiry, err := time.ParseDuration(rc.JWTRefreshExpiry) // refresh token expire: 7d
	if err != nil {
		return nil, fmt.Errorf("config: parse jwt_refresh_expiry %q: %w", rc.JWTRefreshExpiry, err)
	}

	return &Config{
		Port:             rc.Port,
		JWTAccessSecret:  rc.JWTAccessSecret,
		JWTRefreshSecret: rc.JWTRefreshSecret,
		JWTAccessExpiry:  accessExpiry,
		JWTRefreshExpiry: refreshExpiry,
	}, nil
}
