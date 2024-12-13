package config

import "time"

type JWTConfig struct {
	SecretKey       string        `mapstructure:"JWT_SECRET_KEY"`
	AuthTokenTTL    time.Duration `mapstructure:"JWT_AUTH_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `mapstructure:"JWT_REFRESH_TOKEN_TTL"`
}
