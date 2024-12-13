package utils

import (
	"errors"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

type UserType string

const (
	UserEmployee         UserType = "employee"
	UserCustomer         UserType = "customer"
	ACCESS_TOKEN_HEADER           = "Authorization"
	REFRESH_TOKEN_HEADER          = "Refresh-Token"
)

type UserClaims interface {
	jwt.Claims
	GetUserType() UserType
}

type EmployeeClaims struct {
	jwt.RegisteredClaims
	ID           uint              `json:"id"`
	Role         data.EmployeeRole `json:"role"`
	WorkplaceID  uint              `json:"workplaceId"`
	EmployeeType data.EmployeeType `json:"workplaceType"`
}

func (ec *EmployeeClaims) GetUserType() UserType {
	return UserEmployee
}

type CustomerClaims struct {
	jwt.RegisteredClaims
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

func (ec *CustomerClaims) GetUserType() UserType {
	return UserCustomer
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateJWT(claims UserClaims) (string, error) {
	var secretKey string

	cfg := config.GetConfig()

	switch claims.GetUserType() {
	case UserEmployee:
		secretKey = cfg.JWT.EmployeeSecretKey
	case UserCustomer:
		secretKey = cfg.JWT.CustomerSecretKey
	default:
		return "", errors.New("unsupported user type")
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenWithClaims.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
		return errors.New("invalid employee token")
	}

	return nil
}

func ValidateCustomerJWT(tokenString string, claims jwt.Claims) error {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.CustomerSecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid customer token")
	}

	return nil
}

func ExtractEmployeeAccessTokenAndValidate(c *gin.Context) (*EmployeeClaims, error) {
	authHeader := c.GetHeader(ACCESS_TOKEN_HEADER)
	var tokenString string

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		cookie, err := c.Cookie(employees.EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY)

		if err != nil {
			return nil, err
		}
		tokenString = cookie
	}

	claims := &EmployeeClaims{}
	if err := ValidateEmployeeJWT(tokenString, claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func ExtractEmployeeRefreshTokenAndValidate(c *gin.Context) (*EmployeeClaims, error) {
	refreshHeader := c.GetHeader(REFRESH_TOKEN_HEADER)
	var tokenString string

	if refreshHeader != "" && strings.HasPrefix(refreshHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(refreshHeader, "Bearer ")
	} else {
		cookie, err := c.Cookie(employees.EMPLOYEE_REFRESH_TOKEN_COOKIE_KEY)

		if err != nil {
			return nil, err
		}
		tokenString = cookie
	}

	claims := &EmployeeClaims{}
	if err := ValidateEmployeeJWT(tokenString, claims); err != nil {
		return nil, err
	}

	return claims, nil
}
