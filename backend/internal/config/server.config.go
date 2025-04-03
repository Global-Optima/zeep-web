package config

type ServerConfig struct {
	Port              int    `mapstructure:"SERVER_PORT" validate:"required"`
	ClientURL         string `mapstructure:"CLIENT_URL" validate:"required"`
	ImageConverterURL string `mapstructure:"IMAGE_CONVERTER_URL" validate:"required"`
}
