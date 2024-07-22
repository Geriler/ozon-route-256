package consumer_group

import (
	"log"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct{}

func NewConsumerGroupHandler() *ConsumerGroupHandler {
	return &ConsumerGroupHandler{}
}

func (c *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := convertMessage(message)

			log.Printf("%+v\n\n", msg)

			log.Printf(
				"Consumed message\nTopic: %s\nPartition: %d\nOffset: %d\nKey: %s\nPayload: %s\n\n",
				msg.Topic,
				msg.Partition,
				msg.Offset,
				msg.Key,
				msg.Payload,
			)

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

type Msg struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Payload   string `json:"payload"`
}

func convertMessage(in *sarama.ConsumerMessage) Msg {
	return Msg{
		Topic:     in.Topic,
		Partition: in.Partition,
		Offset:    in.Offset,
		Key:       string(in.Key),
		Payload:   string(in.Value),
	}
}
