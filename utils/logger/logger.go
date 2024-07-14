package logger

import (
	"os"

	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/config"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger

	switch cfg.Env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}
	return logger
}
