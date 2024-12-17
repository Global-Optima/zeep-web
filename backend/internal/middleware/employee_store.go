package middleware

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// MatchesStore middleware must be used in case the url path contains parameter store_id
func MatchesStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		storeID, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
		if err != nil {
			zapLogger.Warn(err)
			utils.SendBadRequestError(c, "invalid store ID")
			c.Abort()
			return
		}

		claims, err := contexts.GetEmployeeClaimsFromCtx(c)
		if err != nil {
			zapLogger.Errorf("could not get claims from employee context: %v", err)
			utils.SendErrorWithStatus(c, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		if claims.Role != data.RoleAdmin {
			if claims.EmployeeType != data.StoreEmployeeType || claims.WorkplaceID != uint(storeID) {
				zapLogger.Warnf("employeeID = %d is not assigned to the storeID = %d", claims.EmployeeClaimsData.ID, storeID)
				utils.SendErrorWithStatus(
					c,
					"employee is not assigned to this store",
					http.StatusForbidden)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
