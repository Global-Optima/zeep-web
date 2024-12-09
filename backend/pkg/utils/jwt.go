package utils

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTExpiration = 2 * time.Hour
)

type BaseClaims struct {
	ID   uint   `json:"id"`   // Unified ID for employees/customers
	Type string `json:"type"` // employee(store/warehouse) or customer
	jwt.RegisteredClaims
}

type EmployeeClaims struct {
	BaseClaims
	Role          string `json:"role"`
	WorkplaceID   *uint  `json:"workplace_id,omitempty"`   // Store or Warehouse ID
	WorkplaceType string `json:"workplace_type,omitempty"` // "Store" or "Warehouse"
}

type CustomerClaims struct {
	BaseClaims
	IsVerified bool `json:"is_verified"`
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
