package main

import (
	"trip-advisor/internal/config"
	"trip-advisor/internal/http-server/server"
	"trip-advisor/logger"
	"trip-advisor/router"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.Env)

	logger.Info("Start config", "config", cfg)

	router := router.SetupRouter(cfg)

	server.StartServer(logger, router, cfg)
}
