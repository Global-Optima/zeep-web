package config

import "time"

type JWTConfig struct {
	CustomerSecretKey       string        `mapstructure:"JWT_CUSTOMER_SECRET_KEY" validate:"required"`
	EmployeeSecretKey       string        `mapstructure:"JWT_EMPLOYEE_SECRET_KEY" validate:"required"`
	CustomerAccessTokenTTL  time.Duration `mapstructure:"JWT_CUSTOMER_ACCESS_TOKEN_TTL" default:"1h"`
	CustomerRefreshTokenTTL time.Duration `mapstructure:"JWT_CUSTOMER_REFRESH_TOKEN_TTL" default:"720h"`
	EmployeeAccessTokenTTL  time.Duration `mapstructure:"JWT_EMPLOYEE_ACCESS_TOKEN_TTL" default:"15m"`
	EmployeeRefreshTokenTTL time.Duration `mapstructure:"JWT_EMPLOYEE_REFRESH_TOKEN_TTL" default:"24h"`
}
