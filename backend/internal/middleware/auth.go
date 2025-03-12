package middleware

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	authTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

func EmployeeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, err := authTypes.ExtractEmployeeAccessTokenAndValidate(c)
		if err != nil {
			zapLogger.Warn("missing or invalid token")
			utils.SendErrorWithStatus(c, "missing or invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		contexts.SetEmployeeCtx(c, claims)
		c.Next()
	}
}

func CustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, err := authTypes.ExtractCustomerAccessTokenAndValidate(c)
		if err != nil {
			zapLogger.Warn("missing or invalid token")
			utils.SendErrorWithStatus(c, "missing or invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		contexts.SetCustomerCtx(c, claims)
		c.Next()
	}
}
