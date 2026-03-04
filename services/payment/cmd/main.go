package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shownest/pkg/config"
	"github.com/shownest/pkg/db"
)

func main() {
	ctx := context.Background()

	provider, err := config.NewConfigProvider(ctx)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pool, err := db.New(ctx, provider)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	fmt.Println("connected to postgres successfully")
}
