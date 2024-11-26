package kafka

import (
	"context"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	client sarama.ConsumerGroup
}

func NewKafkaConsumer() (*KafkaConsumer, error) {
	kafkaConfig := config.GetConfig().Kafka

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumerGroup, err := sarama.NewConsumerGroup(kafkaConfig.Brokers, kafkaConfig.ConsumerGroupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{client: consumerGroup}, nil
}

func (c *KafkaConsumer) Consume(topics []string, handler sarama.ConsumerGroupHandler) error {
	ctx := context.Background()

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
