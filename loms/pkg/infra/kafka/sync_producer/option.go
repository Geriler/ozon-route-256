package sync_producer

import (
	"time"

	"github.com/IBM/sarama"
)

type Option func(c *sarama.Config)

func WithProducerPartitioner(pfn sarama.PartitionerConstructor) Option {
	return func(c *sarama.Config) {
		c.Producer.Partitioner = pfn
	}
}

func WithRequiredAcks(acks sarama.RequiredAcks) Option {
	return func(c *sarama.Config) {
		c.Producer.RequiredAcks = acks
	}
}

func WithIdempotent() Option {
	return func(c *sarama.Config) {
		c.Producer.Idempotent = true
	}
}

func WithMaxRetries(n int) Option {
	return func(c *sarama.Config) {
		c.Producer.Retry.Max = n
	}
}

func WithRetryBackoff(d time.Duration) Option {
	return func(c *sarama.Config) {
		c.Producer.Retry.Backoff = d
	}
}

func WithMaxOpenRequests(n int) Option {
	return func(c *sarama.Config) {
		c.Net.MaxOpenRequests = n
	}
}

func WithProducerFlushMessages(n int) Option {
	return func(c *sarama.Config) {
		c.Producer.Flush.Messages = n
	}
}

func WithProducerFlushFrequency(d time.Duration) Option {
	return func(c *sarama.Config) {
		c.Producer.Flush.Frequency = d
	}
}
