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

var logger *slog.Logger
var initError error

func init() {
	cfg, err := config.LoadConfig()
	if err != nil {
		initError = err
		return
	}
	switch cfg.Env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}

func Logger() *slog.Logger {
	return logger
}

func InitError() error {
	return initError
}
