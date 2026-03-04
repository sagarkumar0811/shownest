package config

import (
	"context"
	"os"
)

// NewConfigProvider creates a new ConfigProvider based on the environment variable.
func NewConfigProvider(ctx context.Context) (ConfigProvider, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = EnvironmentLocal
	}

	switch env {
	case EnvironmentProduction, EnvironmentStaging:
		return NewAWSProvider(), nil
	default:
		return NewLocalProvider(DefaultConfigFile)
	}
}
