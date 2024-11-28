package kafka

import (
	"log"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(topic, key string, value []byte) error
	Close() error
}

type KafkaProducer struct {
	client sarama.SyncProducer
}

func NewKafkaProducer(cfg config.KafkaConfig) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = cfg.RetryAttempts
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{client: producer}, nil
}

func (p *KafkaProducer) SendMessage(msg *sarama.ProducerMessage) error {
	partition, offset, err := p.client.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic %s [partition: %d, offset: %d]\n", msg.Topic, partition, offset)
	return nil
}

func (p *KafkaProducer) Close() error {
	err := p.client.Close()
	if err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
	}
	return err
}
