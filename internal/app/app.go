package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"project_sem/internal/config"
	handler "project_sem/internal/handlers"
	"project_sem/internal/repository"
	"project_sem/internal/server"
	"project_sem/internal/service"
	"project_sem/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

func Run(quitSignal <-chan os.Signal) {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config with error %s", err.Error()))
		os.Exit(1)
	}

	postgresStorage, err := storage.NewPostgresPool(ctx, cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create storage %s", err.Error()))
		os.Exit(1)
	}

	migrate, err := storage.NewMigrations(postgresStorage.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to initialize migrations %s", err.Error()))
		os.Exit(1)
	}

	err = migrate.Up()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to apply migrations %s", err.Error()))
		os.Exit(1)
	}

	decimal.MarshalJSONWithoutQuotes = true

	marketingRepository := repository.NewMarketingRepository(postgresStorage)

	marketingService := service.NewMarketingService(marketingRepository)

	marketingApp := fiber.New()

	handler.RegisterRoutes(marketingApp, marketingService)

	quit := make(chan struct{})
	go func() {
		<-quitSignal
		close(quit)
	}()

	marketingDone := make(chan struct{})

	marketingServer := server.New(marketingApp, marketingDone, cfg.APIHost, cfg.GraceTimeout, quit)

	go marketingServer.Run()

	<-marketingDone
}
