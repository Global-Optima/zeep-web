package config

import "time"

type JWTConfig struct {
	CustomerSecretKey       string        `mapstructure:"JWT_CUSTOMER_SECRET_KEY"`
	EmployeeSecretKey       string        `mapstructure:"JWT_EMPLOYEE_SECRET_KEY"`
	CustomerAccessTokenTTL  time.Duration `mapstructure:"JWT_CUSTOMER_ACCESS_TOKEN_TTL"`
	CustomerRefreshTokenTTL time.Duration `mapstructure:"JWT_CUSTOMER_REFRESH_TOKEN_TTL"`
	EmployeeAccessTokenTTL  time.Duration `mapstructure:"JWT_EMPLOYEE_ACCESS_TOKEN_TTL"`
	EmployeeRefreshTokenTTL time.Duration `mapstructure:"JWT_EMPLOYEE_REFRESH_TOKEN_TTL"`
}
