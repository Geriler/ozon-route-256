package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"route256/loms/internal/config"
	"route256/loms/internal/outbox/repository"
	"route256/loms/internal/outbox/service"
)

type Outbox struct {
	cfg     config.Config
	logger  *slog.Logger
	service []*service.OutboxService
}

func NewOutbox(cfg config.Config, logger *slog.Logger) (*Outbox, error) {
	pools, err := dbConnect(context.Background(), cfg.Database.DSNs)
	if err != nil {
		return nil, err
	}

	services := make([]*service.OutboxService, 0, len(pools))

	for _, pool := range pools {
		outboxRepository := repository.NewPostgresOutboxRepository(pool, logger)
		outboxService := service.NewOutboxService(outboxRepository)
		services = append(services, outboxService)
	}

	return &Outbox{
		cfg:     cfg,
		logger:  logger,
		service: services,
	}, nil
}

func (o *Outbox) ClearOutbox(ctx context.Context) error {
	ticker := time.NewTicker(o.cfg.Outbox.ClearTableInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for _, s := range o.service {
				err := s.ClearOutbox(ctx, o.cfg.Outbox.OldDataDuration)
				if err != nil {
					o.logger.Error(fmt.Sprintf("Error clearing outbox: %s", err.Error()))
				}
			}
		}
	}
}
