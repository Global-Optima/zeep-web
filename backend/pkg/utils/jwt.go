package utils

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTExpiration = 2 * time.Hour
)

type BaseClaims struct {
	jwt.RegisteredClaims
}

type EmployeeClaims struct {
	BaseClaims
	Role       types.EmployeeRole `json:"role"`
	StoreID    uint               `json:"storeId"`
	EmployeeID uint               `json:"employeeId"`
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
