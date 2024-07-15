package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"route256/cart/internal/app"
	"route256/cart/internal/config"
	"route256/cart/pkg/lib/logger"
	"route256/cart/pkg/lib/tracing"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	traceProvider := tracing.MustLoadTraceProvider(cfg)

	grpcClient, err := app.NewGRPCClient(cfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	application := app.NewApp(cfg, log, grpcClient)

	go func() {
		log.Info("Starting HTTP application", "port", cfg.HTTP.Port)
		err = application.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Stopping cart service...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.TimeoutStop)
	defer cancel()

	err = grpcClient.Close()
	if err != nil {
		log.Error(err.Error())
	}

	err = application.Shutdown(ctx)
	if err != nil {
		log.Error(err.Error())
	}

	err = traceProvider.Shutdown(ctx)
	if err != nil {
		log.Error(err.Error())
	}
}
