package contexts

import (
	"slices"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/gin-gonic/gin"
)

// GetRegionId returns the retrieved id and HandlerError
func GetRegionId(c *gin.Context) (*uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	var regionID uint
	if !slices.Contains(data.AdminPermissions, claims.Role) {
		if claims.EmployeeType != data.RegionEmployeeType {
			return nil, ErrInvalidEmployeeType
		}
		regionID = claims.WorkplaceID
	} else {
		regionIdStr := c.Query("regionId")
		if regionIdStr == "" {
			return nil, nil
		}
		id, err := strconv.ParseUint(regionIdStr, 10, 64)
		if err != nil {
			return nil, ErrInvalidStoreID
		}
		regionID = uint(id)
	}
	return &regionID, nil
}

func GetRegionIdWithRole(c *gin.Context) (*uint, data.EmployeeRole, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, "", ErrUnauthorizedAccess
	}

	var regionID uint
	if !slices.Contains(data.AdminPermissions, claims.Role) {
		if claims.EmployeeType != data.RegionEmployeeType {
			return nil, "", ErrInvalidEmployeeType
		}
		regionID = claims.WorkplaceID
	} else {
		regionIdStr := c.Query("regionId")
		if regionIdStr == "" {
			return nil, claims.Role, nil
		}
		id, err := strconv.ParseUint(regionIdStr, 10, 64)
		if err != nil {
			return nil, "", ErrInvalidStoreID
		}
		regionID = uint(id)
	}

	return &regionID, claims.Role, nil
}
