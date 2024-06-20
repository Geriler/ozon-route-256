package main

import (
	"fmt"
	"os"

	"route256/cart/internal/app"
	"route256/cart/internal/config"
	"route256/cart/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	grpcClient, err := app.NewGRPCClient(cfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	application := app.NewApp(cfg, log, grpcClient)

	log.Info(fmt.Sprintf("Starting HTTP application on port %d", cfg.HTTP.Port))

	err = application.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
