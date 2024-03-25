package main

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/registration"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/storage"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/app"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/middle"
)

func main() {
	logger := logger.Logger()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}

	sessionStorage := storage.NewSessionStorage()
	useCase := usecase.NewAuthUseCase(sessionStorage)

	authHandler := authorization.NewAuthorizationHandler(useCase)
	regHandler := registration.NewRegistrationHandler(useCase)
	authMiddleware := middle.NewAuthMiddleware(useCase)

	app := app.NewApp(authHandler, regHandler, authMiddleware)

	logger.Info("Start config", "config", cfg)

	router := router.SetupRouter(app)

	if err := server.StartServer(router, cfg); err != nil {
		logger.Error("Failed to start server", "error", err)
		return
	}
}
