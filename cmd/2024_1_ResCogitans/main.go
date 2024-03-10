package main

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	_ "github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/handlers/authorization" // Swagger должен видеть остальные модули, чтобы делать документацию
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/http-server/server"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
)

// @title Your API Title
// @version 1.0
// @description This is a sample API with Swagger documentation.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger := logger.Logger()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}

	logger.Info("Start config", "config", cfg)

	router := router.SetupRouter(cfg)

	if err := server.StartServer(router, cfg); err != nil {
		logger.Error("Failed to start server", "error", err)
		return
	}
}
