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

	awsConfig, mockMode, err := aws.Init(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: aws: %w", err)
	}

	serviceConfig, err := config.Load(ctx, provider)
	if err != nil {
		return fmt.Errorf("wire: service config: %w", err)
	}

	s3Client := aws.NewS3Client(awsConfig, serviceConfig.S3Bucket, mockMode)
	jwtService := jwt.NewService(
		serviceConfig.JWTAccessSecret,
		serviceConfig.JWTRefreshSecret,
		serviceConfig.JWTAccessExpiry,
		serviceConfig.JWTRefreshExpiry,
	)

	repo := repository.New(pool)
	usecase := usecases.New(repo, s3Client)
	handler := handlers.New(usecase)

	return routes.InitRoutes(routes.Config{
		Handler:    handler,
		JWTService: jwtService,
		Port:       serviceConfig.Port,
	})
}
