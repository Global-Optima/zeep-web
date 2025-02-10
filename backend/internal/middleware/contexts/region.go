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
	ErrInvalidRegionID = handlerErrors.NewHandlerError(errors.New("invalid region ID"), http.StatusBadRequest)
	ErrEmptyRegionID   = handlerErrors.NewHandlerError(errors.New("empty region ID"), http.StatusBadRequest)
)

// GetRegionId returns the retrieved id and HandlerError
func GetRegionId(c *gin.Context) (*uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	var regionID uint
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleOwner {
		if claims.EmployeeType != data.RegionEmployeeType {
			return nil, ErrInvalidEmployeeType
		}
		regionID = claims.WorkplaceID
	} else {
		regionIdStr := c.Query("regionId")
		//TODO change
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
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleOwner {
		if claims.EmployeeType != data.RegionEmployeeType {
			return nil, "", ErrInvalidEmployeeType
		}
		regionID = claims.WorkplaceID
	} else {
		regionIdStr := c.Query("regionId")
		//TODO change
		if regionIdStr == "" {
			return nil, claims.Role, nil
		}
		id, err := strconv.ParseUint(regionIdStr, 10, 64)
		if err != nil {
			return nil, "", ErrInvalidStoreID
		}
		regionID = uint(id)
		return nil, claims.Role, nil
	}

	return &regionID, claims.Role, nil
}
