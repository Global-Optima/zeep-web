package middleware

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EmployeeRoleMiddleware(requiredRoles ...data.EmployeeRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, err := contexts.GetEmployeeClaimsFromCtx(c)
		if err != nil {
			zapLogger.Warnf("missing or invalid token: %v", err)
			utils.SendErrorWithStatus(c, "missing or invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if claims.Role == data.RoleAdmin {
			c.Next()
			return
		}

		for _, role := range requiredRoles {
			if claims.Role == role {
				c.Next()
				return
			}
		}

		utils.SendErrorWithStatus(c, "access denied", http.StatusForbidden)
		c.Abort()
	}
}
