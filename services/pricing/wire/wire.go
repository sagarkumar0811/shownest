package wire

import (
	"context"
	"fmt"

	"github.com/shownest/pkg/cache"
	pkgconfig "github.com/shownest/pkg/config"
	"github.com/shownest/pkg/db"
	"github.com/shownest/pkg/jwt"
	"github.com/shownest/pricing-service/internal/client"
	"github.com/shownest/pricing-service/internal/config"
	"github.com/shownest/pricing-service/internal/handlers"
	"github.com/shownest/pricing-service/internal/repository"
	"github.com/shownest/pricing-service/internal/routes"
	"github.com/shownest/pricing-service/internal/usecases"
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

	// Initialize catalog, inventory, and merchant clients
	catalogClient := client.NewCatalogClient(serviceConfig.ExternalService.Catalog)
	inventoryClient := client.NewInventoryClient(serviceConfig.ExternalService.Inventory)
	merchantClient := client.NewMerchantClient(serviceConfig.ExternalService.Merchant)

	// Initialize repository, use cases, and handlers
	repo := repository.New(pool)
	usecase := usecases.New(repo, cacheClient, catalogClient, inventoryClient, merchantClient)
	handler := handlers.New(usecase)

	return routes.InitRoutes(routes.Config{
		Handler:    handler,
		JWTService: jwtService,
		Port:       serviceConfig.Port,
	})
}
