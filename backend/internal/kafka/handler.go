package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler
}

type KafkaHandler struct {
	ProcessMessage func(topic, key, value string) error
}

func NewKafkaHandler(processMessage func(topic, key, value string) error) *MessageHandler {
	return &MessageHandler{ProcessMessage: processMessage}
}

type MessageHandler struct {
	ProcessMessage func(topic, key, value string) error
}

func (h *MessageHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *MessageHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := h.ProcessMessage(message.Topic, string(message.Key), string(message.Value))
		if err != nil {
			fmt.Printf("Error processing Kafka message: %v\n", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
