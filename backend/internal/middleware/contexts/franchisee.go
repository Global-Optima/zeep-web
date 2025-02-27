package contexts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/gin-gonic/gin"
	"slices"
	"strconv"
)

// GetFranchiseeId returns the retrieved id and HandlerError
func GetFranchiseeId(c *gin.Context) (*uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	var franchiseeID uint
	if !slices.Contains(data.AdminPermissions, claims.Role) {
		if claims.EmployeeType != data.FranchiseeEmployeeType {
			return nil, ErrInvalidEmployeeType
		}
		franchiseeID = claims.WorkplaceID
	} else {
		franchiseeIdStr := c.Query("franchiseeId")
		if franchiseeIdStr == "" {
			return nil, nil
		}
		id, err := strconv.ParseUint(franchiseeIdStr, 10, 64)
		if err != nil {
			return nil, ErrInvalidStoreID
		}
		franchiseeID = uint(id)
	}

	return &franchiseeID, nil
}

func GetFranchiseeIdWithRole(c *gin.Context) (*uint, data.EmployeeRole, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, "", ErrUnauthorizedAccess
	}

	var franchiseeID uint
	if !slices.Contains(data.AdminPermissions, claims.Role) {
		if claims.EmployeeType != data.FranchiseeEmployeeType {
			return nil, "", ErrInvalidEmployeeType
		}
		franchiseeID = claims.WorkplaceID
	} else {
		franchiseeIdStr := c.Query("franchiseeId")
		if franchiseeIdStr == "" {
			return nil, claims.Role, nil
		}
		id, err := strconv.ParseUint(franchiseeIdStr, 10, 64)
		if err != nil {
			return nil, "", ErrInvalidStoreID
		}
		franchiseeID = uint(id)
	}

	return &franchiseeID, claims.Role, nil
}
