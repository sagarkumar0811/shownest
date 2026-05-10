package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shownest/pkg/config"
)

type ExternalService struct {
	Inventory string `json:"inventory"`
	Catalog   string `json:"catalog"`
	Merchant  string `json:"merchant"`
}

type Config struct {
	App             string          `json:"app"`
	Port            string          `json:"port"`
	ExternalService ExternalService `json:"externalService"`
}

func Load(ctx context.Context, provider config.ConfigProvider) (*Config, error) {
	raw, err := provider.Get(ctx, config.ServiceConfig)
	if err != nil {
		return nil, fmt.Errorf("config: get service config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse service config: %w", err)
	}

	return &cfg, nil
}
