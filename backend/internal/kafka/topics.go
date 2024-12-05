package kafka

type Topic string

type Topics struct {
	ActiveOrders    Topic
	CompletedOrders Topic
}

func (k *KafkaManager) GetTopic(topic Topic) string {
	return string(topic)
}
