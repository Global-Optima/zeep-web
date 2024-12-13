package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func MapEmployeeToTokenClaims(employee *data.Employee, expTime time.Duration) *utils.EmployeeClaims {
	var workplaceID uint
	var workplaceType data.EmployeeType

	if employee.StoreEmployee != nil {
		workplaceID = employee.StoreEmployee.StoreID
		workplaceType = data.StoreEmployeeType
	} else if employee.WarehouseEmployee != nil {
		workplaceID = employee.WarehouseEmployee.WarehouseID
		workplaceType = data.WarehouseEmployeeType
	}

	claims := utils.EmployeeClaims{
		ID: employee.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:         employee.Role,
		WorkplaceID:  workplaceID,
		EmployeeType: workplaceType,
	}

	return &claims
}
