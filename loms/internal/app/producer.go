package app

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"route256/loms/internal/config"
	"route256/loms/internal/outbox/model"
	outboxRepository "route256/loms/internal/outbox/repository"
	repository "route256/loms/internal/outbox/repository/sqlc"
	"route256/loms/internal/outbox/service"
	"route256/loms/pkg/infra/kafka/sync_producer"
)

type Producer struct {
	outbox   *service.OutboxService
	producer sarama.SyncProducer
	config   config.Config
	logger   *slog.Logger
}

func NewProducer(config config.Config, logger *slog.Logger) (*Producer, error) {
	conn, err := dbConnect(context.Background(), config.Database.DSN)
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

	outboxRepo := outboxRepository.NewPostgresOutboxRepository(conn, logger)
	outbox := service.NewOutboxService(outboxRepo)

	return &Producer{
		outbox:   outbox,
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
			err := p.outbox.SendMessage(ctx, func(ctx context.Context, message *repository.FetchNextMsgsRow) error {
				event := &model.Message{
					OrderID:   message.OrderID,
					EventType: message.EventType,
				}
				bytes, err := json.Marshal(event)
				if err != nil {
					return err
				}

				hasher := md5.New()
				hasher.Write([]byte(fmt.Sprintf("%s-%d", event.EventType, event.OrderID)))
				key := hex.EncodeToString(hasher.Sum(nil))

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
					slog.String("key", key),
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

func (p *Producer) GracefulShutdown() error {
	return p.producer.Close()
}
