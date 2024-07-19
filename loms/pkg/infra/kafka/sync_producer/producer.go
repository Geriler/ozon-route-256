package sync_producer

import (
	"github.com/IBM/sarama"
	"route256/loms/internal/config"
)

func NewSyncProducer(cfg config.Config, opts ...Option) (sarama.SyncProducer, error) {
	syncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Addresses, PrepareConfig(opts...))
	if err != nil {
		return nil, err
	}

	return syncProducer, nil
}
