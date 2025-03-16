package config

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST" validate:"required"`
	Port     int    `mapstructure:"REDIS_PORT" validate:"required"`
	Password string `mapstructure:"REDIS_PASSWORD" validate:"required"`
	Username string `mapstructure:"REDIS_USERNAME" validate:"required"`
	DB       int    `mapstructure:"REDIS_DB" validate:"gte=0"`
}
