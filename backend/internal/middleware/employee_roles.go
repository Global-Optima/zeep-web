package middleware

import (
	"net/http"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func EmployeeRoleMiddleware(requiredRoles ...data.EmployeeRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ExtractEmployeeTokenAndValidate(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			c.Abort()
			return
		}

		for _, role := range requiredRoles {
			if claims.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access"})
		c.Abort()
	}
}

func ExtractEmployeeTokenAndValidate(c *gin.Context) (*utils.EmployeeClaims, error) {
	authHeader := c.GetHeader("Authorization")
	var tokenString string

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {

		cookie, err := c.Cookie(employees.EMPLOYEE_TOKEN_COOKIE_KEY)
		if err != nil {
			return nil, err
		}
		tokenString = cookie
	}

	claims := &utils.EmployeeClaims{}
	if err := utils.ValidateJWT(tokenString, claims); err != nil {
		return nil, err
	}

	return claims, nil
}
