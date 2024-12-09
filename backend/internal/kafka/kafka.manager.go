package kafka

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/Global-Optima/zeep-web/backend/internal/config"
// 	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
// 	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
// 	"github.com/IBM/sarama"
// )

// type KafkaManager struct {
// 	producer sarama.SyncProducer
// 	client   sarama.Client
// 	Topics   Topics
// }

// type Topic string

// type Topics struct {
// 	ActiveOrders    Topic
// 	CompletedOrders Topic
// }

// var log = logger.NewSugar("KAFKA_MANAGER")

// func NewKafkaManager(cfg config.KafkaConfig) (*KafkaManager, error) {
// 	config := sarama.NewConfig()
// 	config.Version = sarama.V2_8_0_0
// 	config.Producer.Return.Errors = true
// 	config.Producer.Return.Successes = true
// 	config.Producer.RequiredAcks = sarama.WaitForAll

// 	client, err := sarama.NewClient(cfg.Brokers, config)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create kafka client: %w", err)
// 	}

// 	producer, err := sarama.NewSyncProducerFromClient(client)
// 	if err != nil {
// 		client.Close()
// 		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
// 	}

// 	return &KafkaManager{
// 		client:   client,
// 		producer: producer,
// 		Topics: Topics{
// 			ActiveOrders:    Topic(cfg.Topics.ActiveOrders),
// 			CompletedOrders: Topic(cfg.Topics.CompletedOrders),
// 		},
// 	}, nil
// }

// func (k *KafkaManager) Close() error {
// 	if err := k.producer.Close(); err != nil {
// 		log.Errorf("failed to close producer: %v", err)
// 	}
// 	if err := k.client.Close(); err != nil {
// 		log.Errorf("failed to close client: %v", err)
// 	}
// 	return nil
// }

// func (k *KafkaManager) GetTopic(t Topic) string {
// 	return string(t)
// }

// func (k *KafkaManager) PublishOrderEvent(topic Topic, storeID uint, event types.OrderEvent) error {
// 	data, err := json.Marshal(event)
// 	if err != nil {
// 		return fmt.Errorf("failed to serialize event: %w", err)
// 	}
// 	msg := &sarama.ProducerMessage{
// 		Topic: k.GetTopic(topic),
// 		Key:   sarama.StringEncoder(fmt.Sprintf("%d", storeID)),
// 		Value: sarama.ByteEncoder(data),
// 	}

// 	_, _, err = k.producer.SendMessage(msg)
// 	if err != nil {
// 		return fmt.Errorf("failed to publish message to topic %s: %w", k.GetTopic(topic), err)
// 	}
// 	return nil
// }

// func (k *KafkaManager) storePartition(topic Topic, storeID uint) (int32, error) {
// 	partitions, err := k.client.Partitions(k.GetTopic(topic))
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get partitions: %w", err)
// 	}
// 	if len(partitions) == 0 {
// 		return 0, fmt.Errorf("no partitions for topic %s", topic)
// 	}
// 	numPartitions := int32(len(partitions))
// 	partition := int32(storeID) % numPartitions
// 	return partition, nil
// }

// func (k *KafkaManager) latestOffset(topic Topic, partition int32) (int64, error) {
// 	latestOffset, err := k.client.GetOffset(k.GetTopic(topic), partition, sarama.OffsetNewest)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get latest offset: %w", err)
// 	}
// 	return latestOffset, nil
// }

// func (k *KafkaManager) FetchOrderEvent(orderID, storeID uint, topic Topic) (*types.OrderEvent, error) {
// 	partition, err := k.storePartition(topic, storeID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// For simplicity, start from oldest. Could optimize by starting from a recent offset.
// 	consumer, err := sarama.NewConsumerFromClient(k.client)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create consumer: %w", err)
// 	}
// 	defer consumer.Close()

// 	pc, err := consumer.ConsumePartition(k.GetTopic(topic), partition, sarama.OffsetOldest)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to consume partition: %w", err)
// 	}
// 	defer pc.Close()

// 	timeout := time.After(3 * time.Second)
// 	for {
// 		select {
// 		case msg := <-pc.Messages():
// 			if msg == nil {
// 				continue
// 			}
// 			var ev types.OrderEvent
// 			if err := json.Unmarshal(msg.Value, &ev); err != nil {
// 				log.Errorf("unmarshal error: %v", err)
// 				continue
// 			}
// 			if ev.StoreID == storeID && ev.ID == orderID {
// 				return &ev, nil
// 			}
// 		case <-timeout:
// 			return nil, fmt.Errorf("timeout while fetching order event %d for store %d", orderID, storeID)
// 		}
// 	}
// }

// func (k *KafkaManager) FetchOrderEvents(topic Topic, storeID uint) ([]types.OrderEvent, error) {
// 	partition, err := k.storePartition(topic, storeID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	latestOffset, err := k.latestOffset(topic, partition)
// 	if err != nil {
// 		return nil, err
// 	}

// 	const maxMessages = 1000
// 	startOffset := latestOffset - maxMessages
// 	if startOffset < 0 {
// 		startOffset = sarama.OffsetOldest
// 	}

// 	consumer, err := sarama.NewConsumerFromClient(k.client)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create consumer: %w", err)
// 	}
// 	defer consumer.Close()

// 	pc, err := consumer.ConsumePartition(k.GetTopic(topic), partition, startOffset)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to consume partition: %w", err)
// 	}
// 	defer pc.Close()

// 	var messages []types.OrderEvent
// 	timeout := time.After(3 * time.Second)
// 	for {
// 		select {
// 		case msg := <-pc.Messages():
// 			if msg == nil {
// 				continue
// 			}
// 			var ev types.OrderEvent
// 			if err := json.Unmarshal(msg.Value, &ev); err != nil {
// 				log.Warnf("unmarshal error: %v", err)
// 				continue
// 			}
// 			messages = append(messages, ev)
// 		case <-timeout:
// 			return messages, nil
// 		}
// 	}
// }

// func (k *KafkaManager) StreamOrderEvents(topic Topic, storeID uint, ctx context.Context) (<-chan types.OrderEvent, error) {
// 	partition, err := k.storePartition(topic, storeID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	latestOffset, err := k.latestOffset(topic, partition)
// 	if err != nil {
// 		return nil, err
// 	}

// 	consumer, err := sarama.NewConsumerFromClient(k.client)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create consumer: %w", err)
// 	}
// 	// Note: we do not defer Close() here, we close when ctx is done.

// 	pc, err := consumer.ConsumePartition(k.GetTopic(topic), partition, latestOffset)
// 	if err != nil {
// 		consumer.Close()
// 		return nil, fmt.Errorf("failed to consume partition: %w", err)
// 	}

// 	eventCh := make(chan types.OrderEvent, 100)
// 	go func() {
// 		defer close(eventCh)
// 		defer pc.Close()
// 		defer consumer.Close()

// 		for {
// 			select {
// 			case msg := <-pc.Messages():
// 				if msg == nil {
// 					continue
// 				}
// 				var ev types.OrderEvent
// 				if err := json.Unmarshal(msg.Value, &ev); err != nil {
// 					log.Warnf("unmarshal error: %v", err)
// 					continue
// 				}
// 				select {
// 				case eventCh <- ev:
// 				case <-ctx.Done():
// 					return
// 				}
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()

// 	return eventCh, nil
// }
