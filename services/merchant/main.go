package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shownest/merchant-service/wire"
	"github.com/shownest/pkg/config"
	"github.com/shownest/pkg/logger"
)

func main() {
	ctx := context.Background()

	provider, err := config.NewConfigProvider(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := logger.Init(ctx, provider); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	if err := wire.InitializeApp(ctx, provider); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
