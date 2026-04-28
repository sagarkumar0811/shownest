package wire

import (
	"context"
	"fmt"

	"github.com/shownest/merchant-service/internal/config"
	"github.com/shownest/merchant-service/internal/handlers"
	"github.com/shownest/merchant-service/internal/repository"
	"github.com/shownest/merchant-service/internal/routes"
	"github.com/shownest/merchant-service/internal/usecases"
	"github.com/shownest/pkg/aws"
	"github.com/shownest/pkg/cache"
	pkgconfig "github.com/shownest/pkg/config"
	"github.com/shownest/pkg/db"
	"github.com/shownest/pkg/jwt"
)

func InitializeApp(ctx context.Context, provider pkgconfig.ConfigProvider) error {

	// Initialize database connection pool
	pool, err := db.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: db: %w", err)
	}
	defer pool.Close()

	// Initialize cache client
	cacheClient, err := cache.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: cache: %w", err)
	}
	defer cacheClient.Close()

	// Load AWS configuration
	awsCfg, cfg, err := aws.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: aws: %w", err)
	}

	// Load service config
	serviceConfig, err := config.Load(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: service config: %w", err)
	}

	// Initialize JWT service
	jwtService, err := jwt.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: jwt: %w", err)
	}

	// Initialize AWS clients and JWT service
	s3Client := aws.NewS3Client(awsCfg, cfg.S3.Bucket, cfg.MockMode)

	// Initialize repository, use cases, and handlers
	repo := repository.New(pool)
	usecase := usecases.New(repo, s3Client, serviceConfig)
	handler := handlers.New(usecase)

	return routes.InitRoutes(routes.Config{
		Handler:    handler,
		JWTService: jwtService,
		Cache:      cacheClient,
		Port:       serviceConfig.Port,
	})
}
