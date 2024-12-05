package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/kafka/setup"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/IBM/sarama"
)

type KafkaManager struct {
	Producer       setup.KafkaProducer
	Consumer       setup.KafkaConsumer
	Topics         Topics
	ConsumeTimeOut time.Duration
}

var Logger = logger.GetInstance()
var DefaultPartition = 0

func NewKafkaManager(cfg config.KafkaConfig) (*KafkaManager, error) {
	producer, err := setup.NewKafkaProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	consumer, err := setup.NewKafkaConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka consumer: %w", err)
	}

	return &KafkaManager{
		Producer: *producer,
		Consumer: *consumer,
		Topics: Topics{
			ActiveOrders:    Topic(cfg.Topics.ActiveOrders),
			CompletedOrders: Topic(cfg.Topics.CompletedOrders),
		},
		ConsumeTimeOut: 3 * time.Second,
	}, nil
}

func (k *KafkaManager) PublishOrderEvent(topic Topic, storeID uint, event types.OrderEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: k.GetTopic(topic),
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", storeID)),
		Value: sarama.ByteEncoder(eventData),
	}

	err = k.Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to publish message to Kafka topic %s: %w", k.GetTopic(topic), err)
	}

	return nil
}

func (k *KafkaManager) FetchOrderEvent(orderID, storeID uint, topic string) (*types.OrderEvent, error) {
	client, err := sarama.NewClient(k.Consumer.Brokers, sarama.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, int32(DefaultPartition), sarama.OffsetOldest)
	if err != nil {
		return nil, fmt.Errorf("failed to consume partition: %w", err)
	}
	defer partitionConsumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), k.ConsumeTimeOut)
	defer cancel()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var event types.OrderEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("Failed to unmarshal Kafka message: %v\n", err)
				continue
			}

			// Match the orderID and storeID
			if event.ID == orderID && event.StoreID == storeID {
				return &event, nil
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("timeout while fetching order event for OrderID %d", orderID)
		}
	}
}

func (k *KafkaManager) FetchOrders(topic string) ([]types.OrderEvent, error) {
	client, err := sarama.NewClient(k.Consumer.Brokers, sarama.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, int32(DefaultPartition), sarama.OffsetOldest)
	if err != nil {
		return nil, fmt.Errorf("failed to consume partition: %w", err)
	}
	defer partitionConsumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), k.ConsumeTimeOut)
	defer cancel()

	var messages []types.OrderEvent
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var event types.OrderEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("Failed to unmarshal Kafka message: %v\n", err)
				continue
			}
			messages = append(messages, event)

		case <-ctx.Done():
			return messages, nil
		}
	}
}
