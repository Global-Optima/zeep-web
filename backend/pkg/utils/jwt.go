package utils

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserType string
type TokenType string

const (
	TokenAccess  TokenType = "access"
	TokenRefresh TokenType = "refresh"

	ACCESS_TOKEN_HEADER  = "Authorization"
	REFRESH_TOKEN_HEADER = "Refresh-Token"

	EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY  = "EMPLOYEE_ACCESS_TOKEN"
	EMPLOYEE_REFRESH_TOKEN_COOKIE_KEY = "EMPLOYEE_REFRESH_TOKEN"
	CUSTOMER_ACCESS_TOKEN_COOKIE_KEY  = "CUSTOMER_ACCESS_TOKEN"
	CUSTOMER_REFRESH_TOKEN_COOKIE_KEY = "CUSTOMER_REFRESH_TOKEN"
)

type EmployeeClaimsData struct {
	ID           uint              `json:"id"`
	Role         data.EmployeeRole `json:"role"`
	WorkplaceID  uint              `json:"workplaceId"`
	EmployeeType data.EmployeeType `json:"workplaceType"`
}

type EmployeeClaims struct {
	jwt.RegisteredClaims
	EmployeeClaimsData
}

type CustomerClaimsData struct {
	ID         uint `json:"id"`
	IsVerified bool `json:"isVerified"`
}

type CustomerClaims struct {
	jwt.RegisteredClaims
	CustomerClaimsData
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateCustomerJWT(input *CustomerClaimsData, tokenType TokenType) (string, error) {
	var ttl time.Duration

	cfg := config.GetConfig()

	switch tokenType {
	case TokenAccess:
		ttl = cfg.JWT.CustomerAccessTokenTTL
	case TokenRefresh:
		ttl = cfg.JWT.CustomerRefreshTokenTTL
	default:
		return "", errors.New("unsupported token type")
	}

	claims := CustomerClaims{
		CustomerClaimsData: *input,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Issuer:    "zeep-web",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenWithClaims.Header["tokenType"] = tokenType

	tokenString, err := tokenWithClaims.SignedString([]byte(cfg.JWT.CustomerSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateEmployeeJWT(input *EmployeeClaimsData, tokenType TokenType) (string, error) {
	var ttl time.Duration

	cfg := config.GetConfig()

	switch tokenType {
	case TokenAccess:
		ttl = cfg.JWT.EmployeeAccessTokenTTL
	case TokenRefresh:
		ttl = cfg.JWT.EmployeeRefreshTokenTTL
	default:
		return "", errors.New("unsupported token type")
	}

	claims := EmployeeClaims{
		EmployeeClaimsData: *input,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Issuer:    "zeep-web",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenWithClaims.Header["tokenType"] = tokenType

	tokenString, err := tokenWithClaims.SignedString([]byte(cfg.JWT.EmployeeSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateCustomerJWT(tokenString string, claims *CustomerClaims, tokenType TokenType) error {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.CustomerSecretKey), nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	err = validateJWT(token, tokenType)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	return nil
}

func ValidateEmployeeJWT(tokenString string, claims *EmployeeClaims, tokenType TokenType) error {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.EmployeeSecretKey), nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	err = validateJWT(token, tokenType)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	return nil
}

func validateJWT(token *jwt.Token, expectedType TokenType) error {
	if !token.Valid {
		return errors.New("token is invalid")
	}

	tokenType, ok := token.Header["tokenType"].(string)
	if !ok || TokenType(tokenType) != expectedType {
		return fmt.Errorf("invalid token type: expected %s, got %s", expectedType, tokenType)
	}

	return nil
}
