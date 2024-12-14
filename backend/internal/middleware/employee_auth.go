package middleware

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EmployeeIdentity() gin.HandlerFunc {

	return func(c *gin.Context) {

		claims, err := ExtractEmployeeTokenAndValidate(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		//TODO auth and refresh tokens logic

		contexts.SetEmployeeCtx(c, claims)
		c.Next()
	}
}
