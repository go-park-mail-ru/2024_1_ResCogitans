package main

import (
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/server"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"

	_ "github.com/go-park-mail-ru/2024_1_ResCogitans/cmd/2024_1_ResCogitans/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server seller server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger := logger.Logger()
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}
	logger.Info("Start config")

	pdb, rdb, err := initialization.DataBaseInitialization()
	if err != nil {
		logger.Error("DataBase initialization error", "error", err)
	}

	storages := initialization.StorageInit(pdb, rdb)
	usecases := initialization.UseCaseInit(storages)
	handlers := initialization.HandlerInit(usecases)

	router := router.SetupRouter(cfg, handlers)

	if err := server.StartServer(router, cfg); err != nil {
		logger.Error("Failed to start server", "error", err)
		return
	}
}
