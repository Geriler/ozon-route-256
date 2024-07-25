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
	service *service.OutboxService
}

func NewOutbox(cfg config.Config, logger *slog.Logger) (*Outbox, error) {
	conn, err := dbConnect(context.Background(), cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	outboxRepository := repository.NewPostgresOutboxRepository(conn, logger)
	outboxService := service.NewOutboxService(outboxRepository)

	return &Outbox{
		cfg:     cfg,
		logger:  logger,
		service: outboxService,
	}, nil
}

func (o *Outbox) ClearOutbox(ctx context.Context) error {
	ticker := time.NewTicker(o.cfg.Outbox.ClearTableInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := o.service.ClearOutbox(ctx, o.cfg.Outbox.OldDataDuration)
			if err != nil {
				o.logger.Error(fmt.Sprintf("Error clearing outbox: %s", err.Error()))
			}
		}
	}
}
