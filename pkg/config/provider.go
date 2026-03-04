package config

import "context"

// ConfigProvider defines the interface for retrieving configuration values.
type ConfigProvider interface {
	Get(ctx context.Context, key string) ([]byte, error)
}
