package middleware

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EmployeeIdentity() gin.HandlerFunc {

	return func(c *gin.Context) {

		claims, err := utils.ExtractEmployeeAccessTokenAndValidate(c)
		if err != nil {
			claims, err = utils.ExtractEmployeeRefreshTokenAndValidate(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
				return
			}
			c.Redirect(http.StatusTemporaryRedirect, "/api/v1/employees/refresh")
		}

		//TODO change all routes connected to auth to /api/v1/auth to avoid circular middleware check???

		contexts.SetEmployeeCtx(c, claims)
		c.Next()
	}
}
