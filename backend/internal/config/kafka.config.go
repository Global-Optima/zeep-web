package config

type KafkaConfig struct {
	Brokers         []string `mapstructure:"KAFKA_BROKERS"`
	ConsumerGroupID string   `mapstructure:"KAFKA_CONSUMER_GROUP_ID"`
	RetryAttempts   int      `mapstructure:"KAFKA_RETRY_ATTEMPTS"`
	RetryInterval   int      `mapstructure:"KAFKA_RETRY_INTERVAL"`
	Topics          struct {
		ActiveOrders string `mapstructure:"KAFKA_TOPIC_ACTIVE_ORDERS"`
	} `mapstructure:",squash"`
}
