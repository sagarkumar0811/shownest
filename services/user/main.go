package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shownest/pkg/config"
	"github.com/shownest/pkg/logger"
	"github.com/shownest/user-service/wire"
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
	if err := logger.Init(ctx, provider); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Wire and start the server
	if err := wire.InitializeApp(ctx, provider); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
