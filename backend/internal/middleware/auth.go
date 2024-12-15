package middleware

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/tokens"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EmployeeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, err := tokens.ExtractEmployeeAccessTokenAndValidate(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		contexts.SetEmployeeCtx(c, claims)
		c.Next()
	}
}

func CustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, err := tokens.ExtractCustomerAccessTokenAndValidate(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		contexts.SetCustomerCtx(c, claims)
		c.Next()
	}
}
