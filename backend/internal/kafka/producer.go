package kafka

import (
	"log"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	client sarama.SyncProducer
}

func NewKafkaProducer() (*KafkaProducer, error) {
	kafkaConfig := config.GetConfig().Kafka

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = kafkaConfig.RetryAttempts
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{client: producer}, nil
}

func (p *KafkaProducer) SendMessage(topic, key, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.client.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic %s [partition: %d, offset: %d]\n", topic, partition, offset)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.client.Close()
}
