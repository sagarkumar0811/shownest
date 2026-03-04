package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type localProvider struct {
	configs map[string]json.RawMessage
}

const configFilePath = "config.json"

// NewLocalProvider returns a new instance of localProvider which implements the ConfigProvider interface.
func NewLocalProvider() (ConfigProvider, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("local config: read file %q: %w", configFilePath, err)
	}

	var configs map[string]json.RawMessage
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("local config: parse %q: %w", configFilePath, err)
	}

	return &localProvider{configs: configs}, nil
}

// Get retrieves the configuration value for the given key from the local config file.
func (p *localProvider) Get(ctx context.Context, key string) ([]byte, error) {
	v, ok := p.configs[key]
	if !ok {
		return nil, fmt.Errorf("local config: key %q not found in config file", key)
	}
	return v, nil
}
