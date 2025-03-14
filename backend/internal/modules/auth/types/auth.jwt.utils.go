package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ACCESS_TOKEN_HEADER = "Authorization"

	EMPLOYEE_SESSION_COOKIE_KEY = "ZEEP_EMPLOYEE_SESSION"
	CUSTOMER_SESSION_COOKIE_KEY = "ZEEP_CUSTOMER_SESSION"
)

type EmployeeClaims struct {
	jwt.RegisteredClaims
	EmployeeID uint `json:"employeeId"`
}

type CustomerClaims struct {
	jwt.RegisteredClaims
	CustomerID uint `json:"customerId"`
}

type Token struct {
	SessionToken string `json:"sessionToken"`
}

func GenerateEmployeeJWT(employeeID uint) (string, error) {
	cfg := config.GetConfig()
	ttl := cfg.JWT.EmployeeTokenTTL

	session := EmployeeClaims{
		EmployeeID: employeeID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Issuer:    "zeep-web",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	return token.SignedString([]byte(cfg.JWT.EmployeeSecretKey))
}

func GenerateCustomerJWT(customerID uint) (string, error) {
	cfg := config.GetConfig()
	ttl := cfg.JWT.EmployeeTokenTTL

	session := CustomerClaims{
		CustomerID: customerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Issuer:    "zeep-web",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	return token.SignedString([]byte(cfg.JWT.EmployeeSecretKey))
}
