package config

type JWTConfig struct {
	SecretKey string `mapstructure:"JWT_SECRET_KEY"`
}
