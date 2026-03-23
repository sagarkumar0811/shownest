package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shownest/catalog-service/internal/routes"
	"github.com/shownest/pkg/cache"
	"github.com/shownest/pkg/config"
	"github.com/shownest/pkg/db"
	"github.com/shownest/pkg/logger"
)

func main() {
	ctx := context.Background()

	// Load configuration
	provider, err := config.NewConfigProvider(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize logger
	err = logger.Init(ctx, provider)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Connect to database
	pool, err := db.Init(ctx, provider)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

	// Connect to cache
	cacheClient, err := cache.Init(ctx, provider)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cacheClient.Close()

	// Start the server
	r := routes.InitRoutes()
	if err := r.Run(":6003"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
