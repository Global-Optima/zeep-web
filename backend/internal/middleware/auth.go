package middleware

import (
	"net/http"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/employeeToken"
	authTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

func EmployeeAuth(employeeTokenManager employeeToken.EmployeeTokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, token, err := authTypes.ExtractEmployeeSessionTokenAndValidate(c)
		if err != nil {
			zapLogger.Warn("missing or invalid token")
			utils.SendErrorWithStatus(c, "missing or invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		savedToken, err := employeeTokenManager.GetTokenByEmployeeID(claims.EmployeeID)
		if err != nil {
			zapLogger.Error("error getting token from db")
			utils.SendErrorWithStatus(c, "error getting token from db", http.StatusInternalServerError)
			c.Abort()
			return
		}

		if savedToken == nil {
			zapLogger.Warn("token not found")
			utils.SendErrorWithStatus(c, "token not found, re-login", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if savedToken.Token != token {
			zapLogger.Warn("token mismatch, re login")
			utils.SendErrorWithStatus(c, "token mismatch", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if savedToken.ExpiresAt.Before(time.Now()) {
			zapLogger.Warn("token expired")
			utils.SendErrorWithStatus(c, "token expired", http.StatusUnauthorized)
			c.Abort()
			return
		}

		employeeSessionData, err := authTypes.MapEmployeeToEmployeeSessionData(&savedToken.Employee)
		if err != nil {
			zapLogger.Error("error mapping employee to employee session data")
			utils.SendErrorWithStatus(c, "error mapping employee to employee session data", http.StatusInternalServerError)
			c.Abort()
			return
		}

		contexts.SetEmployeeCtx(c, employeeSessionData)
		c.Next()
	}
}

func CustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, _, err := authTypes.ExtractCustomerSessionTokenAndValidate(c)
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
