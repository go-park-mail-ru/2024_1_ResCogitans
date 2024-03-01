package server

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"

	"github.com/go-chi/chi/v5"
)

func StartServer(router *chi.Mux, cfg *config.Config) error {
	logger.Logger().Info("Server is starting:", "address", cfg.HTTPServer.Address)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
