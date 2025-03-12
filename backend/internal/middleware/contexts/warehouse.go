package contexts

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	ErrInvalidWarehouseID = handlerErrors.NewHandlerError(errors.New("invalid Warehouse ID"), http.StatusBadRequest)
	ErrEmptyWarehouseID   = handlerErrors.NewHandlerError(errors.New("empty Warehouse ID"), http.StatusBadRequest)

	warehouseExternalRoles = append(data.AdminPermissions, data.RegionPermissions...)
)

// GetWarehouseId returns the retrieved id and HandlerError
func GetWarehouseId(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, ErrUnauthorizedAccess
	}

	var warehouseID uint
	if !slices.Contains(warehouseExternalRoles, claims.Role) {
		if claims.EmployeeType != data.WarehouseEmployeeType {
			return 0, ErrInvalidEmployeeType
		}

		warehouseID = claims.WorkplaceID
	} else {
		warehouseIdStr := c.Query("warehouseId")
		if warehouseIdStr == "" {
			return 0, ErrEmptyWarehouseID
		}

		id, err := strconv.ParseUint(warehouseIdStr, 10, 64)
		if err != nil {
			return 0, ErrInvalidWarehouseID
		}
		warehouseID = uint(id)
	}

	return warehouseID, nil
}

func GetWarehouseIdWithRole(c *gin.Context) (uint, data.EmployeeRole, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, "", ErrUnauthorizedAccess
	}

	var warehouseID uint
	if !slices.Contains(warehouseExternalRoles, claims.Role) {
		if claims.EmployeeType != data.WarehouseEmployeeType {
			return 0, "", ErrInvalidEmployeeType
		}

		warehouseID = claims.WorkplaceID
	} else {
		warehouseIdStr := c.Query("warehouseId")
		if warehouseIdStr == "" {
			return 0, "", ErrEmptyWarehouseID
		}

		id, err := strconv.ParseUint(warehouseIdStr, 10, 64)
		if err != nil {
			return 0, "", ErrInvalidWarehouseID
		}
		warehouseID = uint(id)
	}

	return warehouseID, claims.Role, nil
}
