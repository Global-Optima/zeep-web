package contexts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var (
	ErrInvalidRegionID      = handlerErrors.NewHandlerError(errors.New("invalid region ID"), http.StatusBadRequest)
	ErrEmptyRegionID        = handlerErrors.NewHandlerError(errors.New("empty region ID"), http.StatusBadRequest)
)

// GetRegionId returns the retrieved id and HandlerError
func GetRegionId(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, ErrUnauthorizedAccess
	}

	var regionID uint
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleOwner {
		if claims.EmployeeType != data.WarehouseRegionManagerEmployeeType {
			return 0, ErrInvalidEmployeeType
		}
		regionID = claims.WorkplaceID
	} else {
		regionIdStr := c.Query("regionId")
		if regionIdStr == "" {
			return 0, ErrEmptyRegionID
		}
		id, err := strconv.ParseUint(regionIdStr, 10, 64)
		if err != nil {
			return 0, ErrInvalidRegionID
		}
		regionID = uint(id)
	}

	return regionID, nil
}
