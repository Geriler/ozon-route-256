package main

import (
	"os"
	"os/signal"
	"syscall"

	"route256/loms/internal/app"
	"route256/loms/internal/config"
	"route256/loms/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	grpcApp := app.NewGRPCApp(cfg, log)
	httpgw := app.NewHTTPGW(cfg, log)

	go func() {
		err := grpcApp.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		err := httpgw.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
