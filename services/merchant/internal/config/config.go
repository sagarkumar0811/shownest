package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pkgconfig "github.com/shownest/pkg/config"
)

type Config struct {
	App              string
	Port             string
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
}

type rawConfig struct {
	App              string `json:"app"`
	Port             string `json:"port"`
	JWTAccessSecret  string `json:"jwtAccessSecret"`
	JWTRefreshSecret string `json:"jwtRefreshSecret"`
	JWTAccessExpiry  string `json:"jwtAccessExpiry"`
	JWTRefreshExpiry string `json:"jwtRefreshExpiry"`
}

func Load(ctx context.Context, provider pkgconfig.ConfigProvider) (*Config, error) {
	raw, err := provider.Get(ctx, "service")
	if err != nil {
		return nil, fmt.Errorf("config: get service config: %w", err)
	}

	var rc rawConfig
	if err := json.Unmarshal(raw, &rc); err != nil {
		return nil, fmt.Errorf("config: parse service config: %w", err)
	}

	accessExpiry, err := time.ParseDuration(rc.JWTAccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("config: parse jwt_access_expiry %q: %w", rc.JWTAccessExpiry, err)
	}
	refreshExpiry, err := time.ParseDuration(rc.JWTRefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("config: parse jwt_refresh_expiry %q: %w", rc.JWTRefreshExpiry, err)
	}

	return &Config{
		App:              rc.App,
		Port:             rc.Port,
		JWTAccessSecret:  rc.JWTAccessSecret,
		JWTRefreshSecret: rc.JWTRefreshSecret,
		JWTAccessExpiry:  accessExpiry,
		JWTRefreshExpiry: refreshExpiry,
	}, nil
}
