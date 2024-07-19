package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"route256/loms/internal/app"
	"route256/loms/internal/config"
	"route256/loms/pkg/lib/logger"
	"route256/loms/pkg/lib/tracing"
)

func main() {
	rootCtx := context.Background()

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	traceProvider := tracing.MustLoadTraceProvider(cfg)

	grpcApp, err := app.NewGRPCApp(cfg, log)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	httpgw := app.NewHTTPGW(cfg, log)

	producer, err := app.NewProducer(cfg, log)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	go func() {
		log.Info("Starting gRPC application", "port", cfg.GRPC.Port)
		err := grpcApp.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		log.Info("Starting HTTP application", "port", cfg.HTTP.Port)
		err := httpgw.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		log.Info("Starting producer")
		err := producer.Start(rootCtx)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Stopping loms service...")

	ctx, cancel := context.WithTimeout(rootCtx, cfg.TimeoutStop)
	defer cancel()

	grpcApp.GracefulStop()
	err = httpgw.Shutdown(ctx)
	if err != nil {
		log.Error(err.Error())
	}

	err = traceProvider.Shutdown(ctx)
	if err != nil {
		log.Error(err.Error())
	}

	err = producer.GracefulShutdown()
	if err != nil {
		log.Error(err.Error())
	}
}
