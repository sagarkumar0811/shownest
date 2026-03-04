package config

import "context"

type awsProvider struct{}

// NewAWSProvider returns a new instance of awsProvider which implements the ConfigProvider interface.
func NewAWSProvider() (ConfigProvider, error) {
	return &awsProvider{}, nil
}

// Get retrieves the configuration value for the given key from AWS Secret Manager.
func (p *awsProvider) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
