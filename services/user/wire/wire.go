package wire

import (
	"context"
	"fmt"

	"github.com/shownest/pkg/aws"
	"github.com/shownest/pkg/cache"
	pkgconfig "github.com/shownest/pkg/config"
	"github.com/shownest/pkg/db"
	"github.com/shownest/pkg/jwt"
	"github.com/shownest/user-service/internal/config"
	"github.com/shownest/user-service/internal/handlers"
	"github.com/shownest/user-service/internal/repository"
	"github.com/shownest/user-service/internal/routes"
	"github.com/shownest/user-service/internal/usecases"
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

	// Initialize AWS clients and JWT service
	snsClient := aws.NewSNSClient(awsCfg, cfg.MockMode)
	jwtService := jwt.NewService(
		serviceConfig.JWTAccessSecret,
		serviceConfig.JWTRefreshSecret,
		serviceConfig.JWTAccessExpiry,
		serviceConfig.JWTRefreshExpiry,
	)

	// Initialize repository, use cases, and handlers
	repository := repository.New(pool)
	usecase := usecases.New(repository, cacheClient, snsClient, jwtService)
	handler := handlers.New(usecase)

	return routes.InitRoutes(routes.Config{
		Handler:    handler,
		JWTService: jwtService,
		Port:       serviceConfig.Port,
	})
}
