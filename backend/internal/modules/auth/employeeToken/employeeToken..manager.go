package employeeToken

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type EmployeeTokenManager interface {
	CreateToken(token *data.EmployeeToken) error

	GetTokenByEmployeeID(employeeID uint) (*data.EmployeeToken, error)

	DeleteToken(token *data.EmployeeToken) error

	DeleteTokenByEmployeeID(employeeID uint) error
}

type employeeTokenManager struct {
	db *gorm.DB
}

func NewEmployeeTokenManager(db *gorm.DB) EmployeeTokenManager {
	return &employeeTokenManager{db: db}
}

func (r *employeeTokenManager) CreateToken(token *data.EmployeeToken) error {
	return r.db.Create(token).Error
}

func (r *employeeTokenManager) GetTokenByEmployeeID(employeeID uint) (*data.EmployeeToken, error) {
	var token data.EmployeeToken
	//TODO: add workplace related preloads for employee
	if err := r.db.Model(&data.EmployeeToken{}).
		Preload("Employee").
		Preload("Employee.StoreEmployee").
		Preload("Employee.WarehouseEmployee").
		Preload("Employee.RegionEmployee").
		Preload("Employee.FranchiseeEmployee").
		Preload("Employee.AdminEmployee").
		Where("employee_id = ?", employeeID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *employeeTokenManager) DeleteToken(token *data.EmployeeToken) error {
	return r.db.Unscoped().Delete(token).Error
}

func (r *employeeTokenManager) DeleteTokenByEmployeeID(employeeID uint) error {
	return r.db.Unscoped().Where("employee_id = ?", employeeID).Delete(&data.EmployeeToken{}).Error
}
