package config

type ServerConfig struct {
	Port      int    `mapstructure:"SERVER_PORT" validate:"required"`
	ClientURL string `mapstructure:"CLIENT_URL" validate:"required"`
}
