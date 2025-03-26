package middleware

import (
	"net/http"
	"slices"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

func EmployeeRoleMiddleware(requiredRoles ...data.EmployeeRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, err := contexts.GetEmployeeClaimsFromCtx(c)
		if err != nil {
			zapLogger.Warnf("missing or invalid token: %v", err)
			localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
			c.Abort()
			return
		}

		if claims.Role == data.RoleAdmin {
			c.Next()
			return
		}

		if slices.Contains(requiredRoles, claims.Role) {
			c.Next()
			return
		}

		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		c.Abort()
	}
}
