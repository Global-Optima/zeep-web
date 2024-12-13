package config

import "time"

type JWTConfig struct {
	CustomerSecretKey string        `mapstructure:"JWT_CUSTOMER_SECRET_KEY"`
	EmployeeSecretKey string        `mapstructure:"JWT_EMPLOYEE_SECRET_KEY"`
	AccessTokenTTL    time.Duration `mapstructure:"JWT_ACCESS_TOKEN_TTL"`
	RefreshTokenTTL   time.Duration `mapstructure:"JWT_REFRESH_TOKEN_TTL"`
}
