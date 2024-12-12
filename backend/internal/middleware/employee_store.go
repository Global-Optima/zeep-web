package middleware

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func MatchesStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		storeID, err := strconv.ParseUint(c.Param("store_id"), 10, 64)
		if err != nil {
			utils.SendBadRequestError(c, "invalid store ID")
			c.Abort()
			return
		}

		claims, err := contexts.GetEmployeeClaimsFromCtx(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if claims.Role != data.RoleAdmin {
			if claims.EmployeeType != data.StoreEmployeeType || claims.WorkplaceID != uint(storeID) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "employee is not assigned to this store",
				})
				return
			}
		}

		c.Next()
	}
}
