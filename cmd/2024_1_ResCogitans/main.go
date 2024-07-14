package main

import (
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/server"
	"go.uber.org/dig"
	"golang.org/x/exp/slog"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	container := dig.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logg := logger.NewLogger(cfg)

	if err := container.Provide(func() *config.Config { return cfg }); err != nil {
		log.Fatalf("Failed to provide config: %v", err)
	}
	if err := container.Provide(func() *slog.Logger { return logg }); err != nil {
		log.Fatalf("Failed to provide logger: %v", err)
	}

	if err := registerDependencies(container); err != nil {
		logg.Error("Failed to register dependencies", "error", err)
		return
	}

	err = container.Invoke(func(cfg *config.Config, logger *slog.Logger, r *chi.Mux) {
		logger.Info(fmt.Sprintf("Server is listening on %s", cfg.HTTPServer.Address))

		if err := server.StartServer(r, cfg); err != nil {
			logger.Error("Failed to start server", "error", err)
		}
	})
	if err != nil {
		logg.Error("Failed to start application", "error", err)
	}
}

func registerDependencies(container *dig.Container) error {
	// Создание среза функций-зависимостей, возвращающих ошибку, которая будет в дальнейшем записываться в логи
	providers := []func() error{
		func() error {
			return container.Provide(func(cfg *config.Config, logger *slog.Logger) (*grpc.ClientConn, error) {
				conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.SessionService.Host, cfg.SessionService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
				}
				return conn, nil
			})
		},
		func() error { return container.Provide(initialization.DataBaseInitialization) },
		func() error { return container.Provide(initialization.StorageInit) },
		func() error { return container.Provide(usecase.NewSessionUseCase) },
		func() error { return container.Provide(initialization.UseCaseInit) },
		func() error { return container.Provide(initialization.HandlerInit) },
		func() error { return container.Provide(router.SetupRouter) },
	}

	for _, provider := range providers {
		if err := provider(); err != nil {
			return err
		}
	}
	return nil
}
