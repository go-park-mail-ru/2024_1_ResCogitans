package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/initialization"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/delivery/server"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/router"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/wrapper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

	prometheus.MustRegister(wrapper.RequestsTotal, wrapper.RequestDuration)

	http.Handle("/api/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)
	go http.ListenAndServe(":9090", nil)

	router := router.SetupRouter(cfg, handlers)

	if err := server.StartServer(router, cfg); err != nil {
		logger.Error("Failed to start server", "error", err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Ваш код обработки запроса здесь
	fmt.Fprintf(w, "Hello, World!")
}
