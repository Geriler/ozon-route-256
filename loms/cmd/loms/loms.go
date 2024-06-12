package main

import (
	"route256/loms/internal/config"
	"route256/loms/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	_ = log
}
