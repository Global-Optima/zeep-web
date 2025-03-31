package employeeToken

import (
	"errors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type EmployeeTokenManager interface {
	CreateToken(token *data.EmployeeToken) error
	GetTokenByEmployeeID(employeeID uint) (*data.EmployeeToken, error)
	DeleteToken(token *data.EmployeeToken) error
	DeleteTokenByEmployeeID(employeeID uint) error
	DeleteTokenByStoreEmployeeID(storeEmployeeID uint) error
	DeleteTokenByWarehouseEmployeeID(warehouseEmployeeID uint) error
	DeleteTokenByRegionEmployeeID(regionEmployeeID uint) error
	DeleteTokenByFranchiseeEmployeeID(franchiseeEmployeeID uint) error
	UpdateTokenExpirationByEmployeeID(employeeID uint, newExpiration time.Time) error
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
	err := r.db.Model(&data.EmployeeToken{}).
		Preload("Employee").
		Preload("Employee.StoreEmployee").
		Preload("Employee.WarehouseEmployee").
		Preload("Employee.RegionEmployee").
		Preload("Employee.FranchiseeEmployee").
		Preload("Employee.AdminEmployee").
		Where("employee_id = ?", employeeID).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (r *employeeTokenManager) DeleteToken(token *data.EmployeeToken) error {
	return r.db.Unscoped().Delete(token).Error
}

func (r *employeeTokenManager) DeleteTokenByEmployeeID(employeeID uint) error {
	err := r.db.Unscoped().
		Where("employee_id = ?", employeeID).
		Delete(&data.EmployeeToken{}).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (r *employeeTokenManager) DeleteTokenByStoreEmployeeID(storeEmployeeID uint) error {
	var storeEmp data.StoreEmployee
	if err := r.db.
		Where("id = ?", storeEmployeeID).
		Preload("Employee").
		First(&storeEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(storeEmp.EmployeeID)
}

func (r *employeeTokenManager) DeleteTokenByWarehouseEmployeeID(warehouseEmployeeID uint) error {
	var whEmp data.WarehouseEmployee
	if err := r.db.
		Where("id = ?", warehouseEmployeeID).
		Preload("Employee").
		First(&whEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(whEmp.EmployeeID)
}

func (r *employeeTokenManager) DeleteTokenByRegionEmployeeID(regionEmployeeID uint) error {
	var regionEmp data.RegionEmployee
	if err := r.db.
		Where("id = ?", regionEmployeeID).
		Preload("Employee").
		First(&regionEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(regionEmp.EmployeeID)
}

func (r *employeeTokenManager) DeleteTokenByFranchiseeEmployeeID(franchiseeEmployeeID uint) error {
	var franchiseeEmp data.FranchiseeEmployee
	if err := r.db.
		Where("id = ?", franchiseeEmployeeID).
		Preload("Employee").
		First(&franchiseeEmp).Error; err != nil {
		return err
	}
	return r.DeleteTokenByEmployeeID(franchiseeEmp.EmployeeID)
}

func (r *employeeTokenManager) UpdateTokenExpirationByEmployeeID(employeeID uint, newExpiration time.Time) error {
	return r.db.
		Model(&data.EmployeeToken{}).
		Where("employee_id = ?", employeeID).
		Update("expires_at", newExpiration).
		Error
}
