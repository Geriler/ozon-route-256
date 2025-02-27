package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"route256/loms/internal/config"
	orderModel "route256/loms/internal/order/model"
	"route256/loms/internal/outbox/model"
	outboxRepository "route256/loms/internal/outbox/repository"
	repository "route256/loms/internal/outbox/repository/sqlc"
	"route256/loms/internal/outbox/service"
	"route256/loms/pkg/infra/kafka/sync_producer"
)

type Producer struct {
	outbox   []*service.OutboxService
	producer sarama.SyncProducer
	config   config.Config
	logger   *slog.Logger
}

func NewProducer(config config.Config, logger *slog.Logger) (*Producer, error) {
	pools, err := dbConnect(context.Background(), config.Database.DSNs)
	if err != nil {
		return nil, err
	}

	syncProducer, err := sync_producer.NewSyncProducer(config,
		sync_producer.WithIdempotent(),
		sync_producer.WithRequiredAcks(sarama.WaitForAll),
		sync_producer.WithMaxOpenRequests(1),
		sync_producer.WithMaxRetries(5),
		sync_producer.WithRetryBackoff(10*time.Millisecond),
	)
	if err != nil {
		return nil, err
	}

	outboxes := make([]*service.OutboxService, 0, len(pools))

	for _, pool := range pools {
		outboxRepo := outboxRepository.NewPostgresOutboxRepository(pool, logger)
		outbox := service.NewOutboxService(outboxRepo)
		outboxes = append(outboxes, outbox)
	}

	return &Producer{
		outbox:   outboxes,
		producer: syncProducer,
		config:   config,
		logger:   logger,
	}, nil
}

func (p *Producer) Start(ctx context.Context) error {
	ticker := time.NewTicker(p.config.Kafka.ProducerMessageInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for _, s := range p.outbox {
				err := s.SendMessage(ctx, func(ctx context.Context, message *repository.FetchNextMsgsRow) error {
					event := &model.Message{
						OrderID:   message.OrderID,
						EventType: orderModel.Status(message.EventType),
					}
					bytes, err := json.Marshal(event)
					if err != nil {
						return err
					}

					key := fmt.Sprintf("%d", event.OrderID)

					msg := &sarama.ProducerMessage{
						Topic:     p.config.Kafka.Topic,
						Key:       sarama.StringEncoder(key),
						Value:     sarama.StringEncoder(bytes),
						Timestamp: time.Now(),
					}

					partition, offset, err := p.producer.SendMessage(msg)
					if err != nil {
						return err
					}

					p.logger.Info("Sent message",
						slog.Int("partition", int(partition)),
						slog.Int64("offset", offset),
						slog.String("message", string(bytes)),
						slog.Int("order_id", int(message.OrderID)),
						slog.String("event_type", message.EventType),
					)

					return nil
				})
				if err != nil {
					return err
				}
			}
		}
	}
}

func (p *Producer) GracefulShutdown() error {
	return p.producer.Close()
}
