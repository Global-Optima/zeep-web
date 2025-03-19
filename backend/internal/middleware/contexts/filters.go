package contexts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/gin-gonic/gin"
)

type StoreContextFilter struct {
	StoreID      *uint
	FranchiseeID *uint
}

type WarehouseContextFilter struct {
	WarehouseID *uint
	RegionID    *uint
}

func GetStoreContextFilter(c *gin.Context) (*StoreContextFilter, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	filter := &StoreContextFilter{}

	switch claims.EmployeeType {
	case data.AdminEmployeeType:
		return filter, nil
	case data.FranchiseeEmployeeType:
		filter.FranchiseeID = &claims.WorkplaceID
	case data.StoreEmployeeType:
		filter.StoreID = &claims.WorkplaceID
	default:
		return nil, ErrInvalidEmployeeType
	}

	return filter, nil
}

func GetWarehouseContextFilter(c *gin.Context) (*WarehouseContextFilter, *handlerErrors.HandlerError) {
	claims, err := GetEmployeeClaimsFromCtx(c)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	filter := &WarehouseContextFilter{}

	switch claims.EmployeeType {
	case data.AdminEmployeeType:
		return filter, nil
	case data.RegionEmployeeType:
		filter.RegionID = &claims.WorkplaceID
	case data.WarehouseEmployeeType:
		filter.WarehouseID = &claims.WorkplaceID
	default:
		return nil, ErrInvalidEmployeeType
	}

	return filter, nil
}
