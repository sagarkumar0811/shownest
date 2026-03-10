package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shownest/pkg/cache"
	"github.com/shownest/pkg/config"
	"github.com/shownest/pkg/db"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
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
	pool, err := db.New(ctx, provider)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	logger.Info("connected to postgres successfully")

	// Connect to Redis
	redisClient, err := cache.New(ctx, provider)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer redisClient.Close()

	logger.Info("connected to redis successfully")
	logger.Info("user service started successfully")
}
