package kafka

import (
	"context"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/IBM/sarama"
)

type Consumer interface {
	Consume(ctx context.Context, topics []string, handler ConsumerGroupHandler) error
	Close() error
}

type KafkaConsumer struct {
	client sarama.ConsumerGroup
}

func NewKafkaConsumer(cfg config.KafkaConfig) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{client: consumerGroup}, nil
}

func (c *KafkaConsumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	for {
		err := c.client.Consume(ctx, topics, handler)
		if err != nil {
			return err
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.client.Close()
}
