package middleware

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EmployeeRoleMiddleware(requiredRoles ...data.EmployeeRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := utils.ExtractEmployeeAccessTokenAndValidate(c)
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
