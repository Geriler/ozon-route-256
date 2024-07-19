package main

import (
	"context"
	"os"
	"sync"

	"route256/notifier/internal/config"
	"route256/notifier/pkg/infra/kafka/consumer_group"
	"route256/notifier/pkg/lib/logger"
)

func main() {
	rootCtx := context.Background()

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	//application := app.NewApp(cfg, log)

	var wg = &sync.WaitGroup{}

	handler := consumer_group.NewConsumerGroupHandler()
	cg, err := consumer_group.NewConsumerGroup(
		cfg.Kafka.Addresses,
		"notifier-consumer-group",
		[]string{cfg.Kafka.Topic},
		handler,
	)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	defer cg.Close()

	cg.Run(rootCtx, wg)

	wg.Wait()

	//log.Info("Starting notifier service...")
	//err := application.Run(rootCtx)
	//if err != nil {
	//	log.Error(err.Error())
	//	os.Exit(1)
	//}
}
