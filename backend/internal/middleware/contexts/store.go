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
	ErrInvalidStoreID      = handlerErrors.NewHandlerError(errors.New("invalid store ID"), http.StatusBadRequest)
	ErrEmptyStoreID        = handlerErrors.NewHandlerError(errors.New("empty store ID"), http.StatusBadRequest)
	ErrUnauthorizedAccess  = handlerErrors.NewHandlerError(errors.New("unauthorized access to store"), http.StatusUnauthorized)
	ErrInvalidEmployeeType = handlerErrors.NewHandlerError(errors.New("invalid employee type"), http.StatusBadRequest)

	storeExternalRoles = append(data.AdminPermissions, data.FranchiseePermissions...)
)

// GetStoreId returns the retrieved id and HandlerError
func GetStoreId(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, ErrUnauthorizedAccess
	}

	var storeID uint

	if !slices.Contains(storeExternalRoles, claims.Role) {
		if claims.EmployeeType != data.StoreEmployeeType {
			return 0, ErrInvalidEmployeeType
		}
		storeID = claims.WorkplaceID
	} else {
		storeIdStr := c.Query("storeId")
		if storeIdStr == "" {
			return 0, ErrEmptyStoreID
		}
		id, err := strconv.ParseUint(storeIdStr, 10, 64)
		if err != nil {
			return 0, ErrInvalidStoreID
		}
		storeID = uint(id)
	}

	return storeID, nil
}

func GetStoreIdWithRole(c *gin.Context) (uint, data.EmployeeRole, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, "", ErrUnauthorizedAccess
	}

	var storeID uint
	if !slices.Contains(storeExternalRoles, claims.Role) {
		if claims.EmployeeType != data.StoreEmployeeType {
			return 0, "", ErrInvalidEmployeeType
		}
		storeID = claims.WorkplaceID
	} else {
		storeIdStr := c.Query("storeId")
		if storeIdStr == "" {
			return 0, "", ErrEmptyStoreID
		}
		id, err := strconv.ParseUint(storeIdStr, 10, 64)
		if err != nil {
			return 0, "", ErrInvalidStoreID
		}
		storeID = uint(id)
	}

	return storeID, claims.Role, nil
}
