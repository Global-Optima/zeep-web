package types

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/golang-jwt/jwt/v5"
)

type (
	UserType  string
	TokenType string
)

const (
	TOKEN_TYPE_KEY           = "tokenType"
	TokenAccess    TokenType = "access"
	TokenRefresh   TokenType = "refresh"

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

	tokenWithClaims.Header[TOKEN_TYPE_KEY] = tokenType

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
	tokenWithClaims.Header[TOKEN_TYPE_KEY] = tokenType

	tokenString, err := tokenWithClaims.SignedString([]byte(cfg.JWT.EmployeeSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
