package config

import "time"

type PaymentConfig struct {
	SecretKey   string        `mapstructure:"PAYMENT_SECRET"`
	WaitingTime time.Duration `mapstructure:"PAYMENT_WAIT_TIME"`
}
