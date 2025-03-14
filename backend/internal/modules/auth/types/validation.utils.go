package types

import (
	"fmt"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
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

func ExtractToken(c *gin.Context, headerKey, cookieKey string) (string, error) {
	cookie, err := c.Cookie(cookieKey)
	if err != nil {
		return "", err
	}

	return cookie, nil
}

func ValidateEmployeeToken(tokenString string, session *EmployeeClaims) error {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, session, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.EmployeeSecretKey), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}
	return nil
}

func ValidateCustomerToken(tokenString string, session *CustomerClaims) error {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, session, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.CustomerSecretKey), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}
	return nil
}

func ExtractEmployeeSessionTokenAndValidate(c *gin.Context) (*EmployeeClaims, string, error) {
	return extractAndValidateEmployeeToken(
		c,
	)
}

func ExtractCustomerSessionTokenAndValidate(c *gin.Context) (*CustomerClaims, string, error) {
	return extractAndValidateCustomerToken(
		c,
	)
}

func extractAndValidateEmployeeToken(c *gin.Context) (*EmployeeClaims, string, error) {
	tokenString, err := ExtractToken(c, ACCESS_TOKEN_HEADER, EMPLOYEE_SESSION_COOKIE_KEY)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to extract token")
	}

	claims := &EmployeeClaims{}
	if err := ValidateEmployeeToken(tokenString, claims); err != nil {
		return nil, "", errors.Wrap(err, "failed to validate token")
	}

	return claims, tokenString, nil
}

func extractAndValidateCustomerToken(c *gin.Context) (*CustomerClaims, string, error) {
	tokenString, err := ExtractToken(c, ACCESS_TOKEN_HEADER, CUSTOMER_SESSION_COOKIE_KEY)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to extract token")
	}

	claims := &CustomerClaims{}
	if err := ValidateCustomerToken(tokenString, claims); err != nil {
		return nil, "", errors.Wrap(err, "failed to validate token")
	}

	return claims, tokenString, nil
}
