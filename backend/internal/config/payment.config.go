package config

import "time"

type PaymentConfig struct {
	SecretKey   string        `mapstructure:"PAYMENT_SECRET" validate:"required"`
	WaitingTime time.Duration `mapstructure:"PAYMENT_WAIT_TIME" default:"3m"`
}
