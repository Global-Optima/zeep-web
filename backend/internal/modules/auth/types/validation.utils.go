package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"strings"
)

func ValidateCustomer(input CustomerRegisterDTO) error {
	if strings.TrimSpace(input.FirstName) == "" {
		return moduleErrors.ErrValidation.WithDetails("firstName", "customer first name cannot contain empty values")
	}

	if strings.TrimSpace(input.LastName) == "" {
		return moduleErrors.ErrValidation.WithDetails("lastName", "customer last name cannot contain empty values")
	}

	if !utils.IsValidPhone(input.Phone, "") {
		return errors.New("invalid email format")
	}

	if err := utils.IsValidPassword(input.Password); err != nil {
		return moduleErrors.ErrValidation.WithDetails("password", fmt.Sprintf("password validation failed: %v", err))
	}

	return nil
}

func extractAndValidateEmployeeToken(c *gin.Context, tokenType TokenType) (*EmployeeClaims, error) {
	var cookieKey, headerKey string

	switch tokenType {
	case TokenAccess:
		cookieKey = EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY
		headerKey = ACCESS_TOKEN_HEADER
	case TokenRefresh:
		cookieKey = EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY
		headerKey = REFRESH_TOKEN_HEADER
	default:
		return nil, fmt.Errorf("invalid token type: %v", tokenType)
	}

	tokenString, err := ExtractToken(c, headerKey, cookieKey)
	if err != nil {
		return nil, err
	}

	claims := &EmployeeClaims{}

	if err := ValidateEmployeeJWT(tokenString, claims, tokenType); err != nil {
		return nil, err
	}

	return claims, nil
}

func extractAndValidateCustomerToken(c *gin.Context, tokenType TokenType) (*CustomerClaims, error) {
	var cookieKey, headerKey string

	switch tokenType {
	case TokenAccess:
		cookieKey = CUSTOMER_ACCESS_TOKEN_COOKIE_KEY
		headerKey = ACCESS_TOKEN_HEADER
	case TokenRefresh:
		cookieKey = CUSTOMER_ACCESS_TOKEN_COOKIE_KEY
		headerKey = REFRESH_TOKEN_HEADER
	default:
		return nil, fmt.Errorf("invalid token type: %v", tokenType)
	}

	tokenString, err := ExtractToken(c, headerKey, cookieKey)
	if err != nil {
		return nil, err
	}

	claims := &CustomerClaims{}

	if err := ValidateCustomerJWT(tokenString, claims, tokenType); err != nil {
		return nil, err
	}

	return claims, nil
}

// ExtractToken tries to get token from HTTP header or browser cookie. Header is prioritized.
func ExtractToken(c *gin.Context, headerKey, cookieKey string) (string, error) {
	authHeader := c.GetHeader(headerKey)
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	cookie, err := c.Cookie(cookieKey)
	if err != nil {
		return "", err
	}

	return cookie, nil
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

	tokenType, ok := token.Header[TOKEN_TYPE_KEY].(string)
	if !ok || TokenType(tokenType) != expectedType {
		return fmt.Errorf("invalid token type: expected %s, got %s", expectedType, tokenType)
	}

	return nil
}

func ExtractEmployeeAccessTokenAndValidate(c *gin.Context) (*EmployeeClaims, error) {
	return extractAndValidateEmployeeToken(
		c,
		TokenAccess,
	)
}

func ExtractCustomerAccessTokenAndValidate(c *gin.Context) (*CustomerClaims, error) {
	return extractAndValidateCustomerToken(
		c,
		TokenAccess,
	)
}
