package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/IBM/sarama"
	"github.com/cenkalti/backoff/v4"
)

type KafkaManager struct {
	Producer KafkaProducer
	Consumer KafkaConsumer
	Topics   Topics
}

func NewKafkaManager(cfg config.KafkaConfig) (*KafkaManager, error) {
	producer, err := NewKafkaProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	consumer, err := NewKafkaConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka consumer: %w", err)
	}

	return &KafkaManager{
		Producer: *producer,
		Consumer: *consumer,
		Topics: Topics{
			ActiveOrders:    cfg.Topics.ActiveOrders,
			CompletedOrders: cfg.Topics.CompletedOrders,
		},
	}, nil
}

func (k *KafkaManager) PublishOrderEvent(topic string, storeID uint, event types.OrderEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", storeID)), // Partition by only StoreID
		Value: sarama.ByteEncoder(eventData),
	}

	err = k.Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to publish message to Kafka topic %s: %w", topic, err)
	}

	return nil
}

func (k *KafkaManager) GetOrderEvent(topic string, orderID uint, storeID uint) (*types.OrderEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultChan := make(chan *types.OrderEvent, 1)
	handler := NewKafkaHandler(func(_, key, value string) error {
		var event types.OrderEvent
		if err := json.Unmarshal([]byte(value), &event); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		var retrievedStoreID uint
		_, err := fmt.Sscanf(key, "%d", &retrievedStoreID)
		if err != nil || retrievedStoreID != storeID {
			return nil // Ignore if storeID doesn't match
		}

		if event.ID == orderID {
			select {
			case resultChan <- &event:
			default:
			}
			cancel() // Signal completion
		}
		return nil
	})

	operation := func() error {
		if err := k.Consumer.Consume(ctx, []string{topic}, handler); err != nil {
			return err
		}
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	err := backoff.Retry(operation, expBackoff)
	if err != nil {
		return nil, fmt.Errorf("failed to consume Kafka messages: %w", err)
	}

	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout fetching order event for order ID %d", orderID)
	}
}
