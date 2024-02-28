package server

import (
	"net/http"
	"trip-advisor/internal/config"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

func StartServer(logger *slog.Logger, router *chi.Mux, cfg *config.Config) {
	logger.Info("Server is starting:", "address", cfg.HTTPServer.Address)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Idle_timeoute,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start")
	}
}
