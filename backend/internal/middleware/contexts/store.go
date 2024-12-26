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
	ErrInvalidStoreID      = handlerErrors.NewHandlerError(errors.New("invalid store ID"), http.StatusBadRequest)
	ErrUnauthorizedAccess  = handlerErrors.NewHandlerError(errors.New("unauthorized access to store"), http.StatusUnauthorized)
	ErrInvalidEmployeeType = handlerErrors.NewHandlerError(errors.New("invalid employee type"), http.StatusBadRequest)
)

// GetStoreId returns the retrieved id and HandlerError
func GetStoreId(c *gin.Context) (uint, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return 0, ErrUnauthorizedAccess
	}

	var storeID uint
	if claims.Role != data.RoleAdmin && claims.Role != data.RoleDirector {
		storeID = claims.WorkplaceID
	} else {
		if claims.EmployeeType != data.StoreEmployeeType {
			return 0, ErrInvalidEmployeeType
		}
		
		id, err := strconv.ParseUint(c.Query("storeId"), 10, 64)
		if err != nil {
			return 0, ErrInvalidStoreID
		}
		storeID = uint(id)
	}

	return storeID, nil
}
