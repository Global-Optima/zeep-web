package config

import "time"

const (
	DEFAULT_CUSTOMER_ACCESS_TOKEN_TTL  = 1 * time.Hour
	DEFAULT_CUSTOMER_REFRESH_TOKEN_TTL = 30 * 24 * time.Hour
	DEFAULT_EMPLOYEE_ACCESS_TOKEN_TTL  = 15 * time.Minute
	DEFAULT_EMPLOYEE_REFRESH_TOKEN_TTL = 24 * time.Hour
)

type JWTConfig struct {
	CustomerSecretKey       string        `mapstructure:"JWT_CUSTOMER_SECRET_KEY"`
	EmployeeSecretKey       string        `mapstructure:"JWT_EMPLOYEE_SECRET_KEY"`
	CustomerAccessTokenTTL  time.Duration `mapstructure:"JWT_CUSTOMER_ACCESS_TOKEN_TTL"`
	CustomerRefreshTokenTTL time.Duration `mapstructure:"JWT_CUSTOMER_REFRESH_TOKEN_TTL"`
	EmployeeAccessTokenTTL  time.Duration `mapstructure:"JWT_EMPLOYEE_ACCESS_TOKEN_TTL"`
	EmployeeRefreshTokenTTL time.Duration `mapstructure:"JWT_EMPLOYEE_REFRESH_TOKEN_TTL"`
}
