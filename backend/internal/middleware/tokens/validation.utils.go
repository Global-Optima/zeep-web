package tokens

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func extractAndValidateEmployeeToken(c *gin.Context, tokenType utils.TokenType) (*utils.EmployeeClaims, error) {
	var cookieKey, headerKey string

	switch tokenType {
	case utils.TokenAccess:
		cookieKey = utils.EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY
		headerKey = utils.ACCESS_TOKEN_HEADER
	case utils.TokenRefresh:
		cookieKey = utils.EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY
		headerKey = utils.REFRESH_TOKEN_HEADER
	default:
		return nil, fmt.Errorf("invalid token type: %v", tokenType)
	}

	tokenString, err := extractToken(c, headerKey, cookieKey)
	if err != nil {
		return nil, err
	}

	claims := &utils.EmployeeClaims{}

	if err := utils.ValidateEmployeeJWT(tokenString, claims, tokenType); err != nil {
		return nil, err
	}

	return claims, nil
}

func extractAndValidateCustomerToken(c *gin.Context, tokenType utils.TokenType) (*utils.CustomerClaims, error) {
	var cookieKey, headerKey string

	switch tokenType {
	case utils.TokenAccess:
		cookieKey = utils.CUSTOMER_ACCESS_TOKEN_COOKIE_KEY
		headerKey = utils.ACCESS_TOKEN_HEADER
	case utils.TokenRefresh:
		cookieKey = utils.CUSTOMER_ACCESS_TOKEN_COOKIE_KEY
		headerKey = utils.REFRESH_TOKEN_HEADER
	default:
		return nil, fmt.Errorf("invalid token type: %v", tokenType)
	}

	tokenString, err := extractToken(c, headerKey, cookieKey)
	if err != nil {
		return nil, err
	}

	claims := &utils.CustomerClaims{}

	if err := utils.ValidateCustomerJWT(tokenString, claims, tokenType); err != nil {
		return nil, err
	}

	return claims, nil
}

func extractToken(c *gin.Context, headerKey, cookieKey string) (string, error) {
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

func ExtractEmployeeAccessTokenAndValidate(c *gin.Context) (*utils.EmployeeClaims, error) {
	return extractAndValidateEmployeeToken(
		c,
		utils.TokenAccess,
	)
}

func ExtractEmployeeRefreshTokenAndValidate(c *gin.Context) (*utils.EmployeeClaims, error) {
	return extractAndValidateEmployeeToken(
		c,
		utils.TokenRefresh,
	)
}

func ExtractCustomerAccessTokenAndValidate(c *gin.Context) (*utils.CustomerClaims, error) {
	return extractAndValidateCustomerToken(
		c,
		utils.TokenAccess,
	)
}

func ExtractCustomerRefreshTokenAndValidate(c *gin.Context) (*utils.CustomerClaims, error) {
	return extractAndValidateCustomerToken(
		c,
		utils.TokenRefresh,
	)
}
