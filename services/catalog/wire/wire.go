package wire

import (
	"context"
	"fmt"

	"github.com/shownest/catalog-service/internal/config"
	"github.com/shownest/catalog-service/internal/handlers"
	"github.com/shownest/catalog-service/internal/repository"
	"github.com/shownest/catalog-service/internal/routes"
	"github.com/shownest/catalog-service/internal/usecases"
	"github.com/shownest/pkg/aws"
	"github.com/shownest/pkg/cache"
	pkgconfig "github.com/shownest/pkg/config"
	"github.com/shownest/pkg/db"
	"github.com/shownest/pkg/jwt"
)

func InitializeApp(ctx context.Context, provider pkgconfig.ConfigProvider) error {
	pool, err := db.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: db: %w", err)
	}
	defer pool.Close()

	cacheClient, err := cache.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: cache: %w", err)
	}
	defer cacheClient.Close()

	awsCfg, awsExtCfg, err := aws.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: aws: %w", err)
	}

	serviceConfig, err := config.Load(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: service config: %w", err)
	}

	jwtService, err := jwt.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: jwt: %w", err)
	}

	s3Client := aws.NewS3Client(awsCfg, awsExtCfg.S3.Bucket, awsExtCfg.MockMode)

	repo := repository.New(pool)
	usecase := usecases.New(repo, s3Client, cacheClient, serviceConfig)
	handler := handlers.New(usecase)

	return routes.InitRoutes(routes.Config{
		Handler:    handler,
		JWTService: jwtService,
		Cache:      cacheClient,
		Port:       serviceConfig.Port,
	})
}
