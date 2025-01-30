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
	ErrInvalidFranchiseeID = handlerErrors.NewHandlerError(errors.New("invalid franchisee ID"), http.StatusBadRequest)
	ErrEmptyFranchiseeID   = handlerErrors.NewHandlerError(errors.New("empty franchisee ID"), http.StatusBadRequest)
)

// GetFranchiseeId returns the retrieved id and HandlerError
func GetFranchiseeId(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, ErrUnauthorizedAccess
	}

	var franchiseeID uint
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleOwner {
		if claims.EmployeeType != data.FranchiseeEmployeeType {
			return 0, ErrInvalidEmployeeType
		}
		franchiseeID = claims.WorkplaceID
	} else {
		franchiseeIdStr := c.Query("franchiseeId")
		if franchiseeIdStr == "" {
			return 0, ErrEmptyFranchiseeID
		}
		id, err := strconv.ParseUint(franchiseeIdStr, 10, 64)
		if err != nil {
			return 0, ErrInvalidFranchiseeID
		}
		franchiseeID = uint(id)
	}

	return franchiseeID, nil
}

func GetFranchiseeIdWithRole(c *gin.Context) (uint, data.EmployeeRole, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, "", ErrUnauthorizedAccess
	}

	var franchiseeID uint
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleOwner {
		if claims.EmployeeType != data.FranchiseeEmployeeType {
			return 0, "", ErrInvalidEmployeeType
		}
		franchiseeID = claims.WorkplaceID
	} else {
		franchiseeIdStr := c.Query("franchiseeId")
		if franchiseeIdStr == "" {
			return 0, "", ErrEmptyFranchiseeID
		}
		id, err := strconv.ParseUint(franchiseeIdStr, 10, 64)
		if err != nil {
			return 0, "", ErrInvalidFranchiseeID
		}
		franchiseeID = uint(id)
	}

	return franchiseeID, claims.Role, nil
}
