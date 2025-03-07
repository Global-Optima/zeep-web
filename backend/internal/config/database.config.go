package config

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST" validate:"required"`
	Port     int    `mapstructure:"DB_PORT" validate:"required"`
	User     string `mapstructure:"DB_USER" validate:"required"`
	Password string `mapstructure:"DB_PASSWORD" validate:"required"`
	Name     string `mapstructure:"DB_NAME" validate:"required"`
}
