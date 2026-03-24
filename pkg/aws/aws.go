package aws

import (
	"context"
	"encoding/json"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/shownest/pkg/config"
)

type Config struct {
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	MockMode        bool   `json:"mockMode"`
}

func Init(ctx context.Context, provider config.ConfigProvider) (awssdk.Config, bool, error) {
	raw, err := provider.Get(ctx, config.AWSCredentials)
	if err != nil {
		return awssdk.Config{}, false, fmt.Errorf("aws: get config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return awssdk.Config{}, false, fmt.Errorf("aws: parse config: %w", err)
	}

	if cfg.MockMode {
		return awssdk.Config{}, true, nil
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return awssdk.Config{}, false, fmt.Errorf("aws: load config: %w", err)
	}

	return awsCfg, false, nil
}
