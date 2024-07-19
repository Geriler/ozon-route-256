package consumer_group

import "github.com/IBM/sarama"

type Option func(*sarama.Config)

func WithOffsetsInitial(v int64) Option {
	return func(c *sarama.Config) {
		c.Consumer.Offsets.Initial = v
	}
}

func WithReturnSuccessesEnabled(isEnabled bool) Option {
	return func(c *sarama.Config) {
		c.Producer.Return.Successes = isEnabled
	}
}
