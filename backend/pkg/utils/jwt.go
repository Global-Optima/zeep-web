package utils

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/golang-jwt/jwt/v5"
)

type BaseClaims struct {
	ID   uint   `json:"id"`   // Unified ID for employees/customers
	Type string `json:"type"` // employee(store/warehouse) or customer
}

type EmployeeClaims struct {
	jwt.RegisteredClaims
	ID           uint              `json:"id"`
	Role         data.EmployeeRole `json:"role"`
	WorkplaceID  uint              `json:"workplaceId"`
	EmployeeType data.EmployeeType `json:"workplaceType"`
}

type CustomerClaims struct {
	jwt.RegisteredClaims
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

func GenerateJWT(claims jwt.Claims, expiration time.Duration) (string, error) {
	cfg := config.GetConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.SecretKey))
}

func ValidateJWT(tokenString string, claims jwt.Claims) error {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
