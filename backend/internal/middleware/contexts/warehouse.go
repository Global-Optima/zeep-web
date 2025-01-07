package contexts

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	ErrInvalidWarehouseID = handlerErrors.NewHandlerError(errors.New("invalid warehouse ID"), http.StatusBadRequest)
)

// GetWarehouseId returns the retrieved id and HandlerError
func GetWarehouseId(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, ErrUnauthorizedAccess
	}

	var warehouseID uint
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleDirector {
		warehouseID = claims.WorkplaceID
	} else {
		if claims.EmployeeType != data.WarehouseEmployeeType {
			return 0, ErrInvalidEmployeeType
		}

		id, err := strconv.ParseUint(c.Query("warehouseId"), 10, 64)
		if err != nil {
			return 0, ErrInvalidWarehouseID
		}
		warehouseID = uint(id)
	}

	return warehouseID, nil
}
