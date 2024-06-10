package main

import (
	"os"

	"route256/cart/internal/app"
	"route256/cart/internal/config"
	"route256/cart/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	application := app.NewApp(cfg, log)
	err := application.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
