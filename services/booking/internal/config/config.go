package config

import (
	"context"
	"encoding/json"
	"fmt"

	pkgconfig "github.com/shownest/pkg/config"
)

type ExternalService struct {
	Inventory string `json:"inventory"`
}

type Config struct {
	App             string          `json:"app"`
	Port            string          `json:"port"`
	ExternalService ExternalService `json:"externalService"`
}

func Load(ctx context.Context, provider pkgconfig.ConfigProvider) (*Config, error) {
	raw, err := provider.Get(ctx, pkgconfig.ServiceConfig)
	if err != nil {
		return nil, fmt.Errorf("config: get service config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse service config: %w", err)
	}

	return &cfg, nil
}
