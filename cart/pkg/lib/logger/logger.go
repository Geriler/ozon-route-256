package logger

import (
	"log/slog"
	"os"
)

type Env string

const (
	Local Env = "local"
	Dev   Env = "dev"
	Prod  Env = "prod"
)

func SetupLogger(env Env) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case Local:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Dev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Prod:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
