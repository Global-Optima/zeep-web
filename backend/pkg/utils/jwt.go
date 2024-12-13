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

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateJWT(claims jwt.Claims, expiration time.Duration) (TokenPair, error) {
	cfg := config.GetConfig()

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := tokenWithClaims.SignedString([]byte(cfg.JWT.EmployeeSecretKey))
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, err := tokenWithClaims.SignedString([]byte(cfg.JWT.EmployeeSecretKey))
	if err != nil {
		return TokenPair{}, err
	}

	tokenPair := TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}

func ValidateEmployeeJWT(tokenString string, claims jwt.Claims) error {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.EmployeeSecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
