package config

import "time"

type JWTConfig struct {
	CustomerSecretKey string        `mapstructure:"JWT_CUSTOMER_SECRET_KEY" validate:"required"`
	EmployeeSecretKey string        `mapstructure:"JWT_EMPLOYEE_SECRET_KEY" validate:"required"`
	CustomerTokenTTL  time.Duration `mapstructure:"JWT_CUSTOMER_TOKEN_TTL" default:"168h"`
	EmployeeTokenTTL  time.Duration `mapstructure:"JWT_EMPLOYEE_TOKEN_TTL" default:"168h"`
}
