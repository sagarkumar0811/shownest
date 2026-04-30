package wire

import (
	"context"
	"fmt"

	"github.com/shownest/booking-service/internal/client"
	"github.com/shownest/booking-service/internal/config"
	"github.com/shownest/booking-service/internal/handlers"
	"github.com/shownest/booking-service/internal/repository"
	"github.com/shownest/booking-service/internal/routes"
	"github.com/shownest/booking-service/internal/usecases"
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

	// Initialize inventory client
	inventoryClient := client.NewInventoryClient(serviceConfig.ExternalService.Inventory)

	// Initialize catalog client
	catalogClient := client.NewCatalogClient(serviceConfig.ExternalService.Catalog)

	// Initialize repository, use cases, and handlers
	repo := repository.New(pool)
	usecase := usecases.New(repo, inventoryClient, catalogClient)
	handler := handlers.New(usecase)

	return routes.InitRoutes(routes.Config{
		Handler:    handler,
		JWTService: jwtService,
		Port:       serviceConfig.Port,
	})
}
