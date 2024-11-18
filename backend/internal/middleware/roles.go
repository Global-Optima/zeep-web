package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user role set in the context after authentication
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found in context"})
			c.Abort()
			return
		}

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access"})
			c.Abort()
			return
		}

		c.Next()
	}
}
