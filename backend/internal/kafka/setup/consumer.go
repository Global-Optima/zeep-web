package setup

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/IBM/sarama"
)

type Consumer interface {
	Consume(ctx context.Context, topics []string, handler ConsumerGroupHandler) error
	Close() error
}

type KafkaConsumer struct {
	client  sarama.ConsumerGroup
	Brokers []string
}

func NewKafkaConsumer(cfg config.KafkaConfig) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.ClientID = "zeep-kafka-consumer"
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{client: consumerGroup, Brokers: cfg.Brokers}, nil
}

func (c *KafkaConsumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err() // Exit if the context is canceled
		default:
			err := c.client.Consume(ctx, topics, handler)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err() // Return if the context is canceled
				}
				fmt.Printf("Kafka consumer error: %v. Retrying...\n", err)
				time.Sleep(1 * time.Second) // Retry with a small delay
			}
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.client.Close()
}
