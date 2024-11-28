package kafka

import (
	"context"
	"log"

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

func parseRebalanceStrategy(strategy string) sarama.BalanceStrategy {
	switch strategy {
	case "range":
		return sarama.NewBalanceStrategyRange()
	case "round_robin":
		return sarama.NewBalanceStrategyRoundRobin()
	case "sticky":
		return sarama.NewBalanceStrategySticky()
	default:
		log.Printf("Unknown rebalance strategy '%s', defaulting to round_robin", strategy)
		return sarama.BalanceStrategyRoundRobin
	}
}

func parseOffsetInitial(initial string) int64 {
	switch initial {
	case "newest":
		return sarama.OffsetNewest
	case "oldest":
		return sarama.OffsetOldest
	default:
		log.Printf("Unknown offset initial '%s', defaulting to oldest", initial)
		return sarama.OffsetOldest
	}
}
