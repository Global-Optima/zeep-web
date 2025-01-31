package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StoreEmployeeRepository interface {
	GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]data.StoreEmployee, error)
	GetStoreEmployeeByID(id, storeID uint) (*data.StoreEmployee, error)
	UpdateStoreEmployee(id uint, storeID uint, updateModels *types.UpdateModels) error
}

type storeEmployeeRepository struct {
	db *gorm.DB
}

func NewStoreEmployeeRepository(db *gorm.DB) StoreEmployeeRepository {
	return &storeEmployeeRepository{db: db}
}

func (r *storeEmployeeRepository) GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]data.StoreEmployee, error) {
	var employees []data.StoreEmployee
	query := r.db.Model(&data.StoreEmployee{}).
		Where("store_id = ?", storeID).Preload("Employee")

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where("role = ?", *filter.Role)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR phone ILIKE ? OR email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employees: %w", err)
	}
	return employees, nil
}

func (r *storeEmployeeRepository) GetStoreEmployeeByID(id, storeID uint) (*data.StoreEmployee, error) {
	var storeEmployee data.StoreEmployee
	err := r.db.Model(&data.StoreEmployee{}).
		Preload("Employee").
		Where("id = ? AND store_id = ?", id, storeID).
		First(&storeEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employee by ID: %w", err)
	}
	return &storeEmployee, nil
}

func (r *storeEmployeeRepository) UpdateStoreEmployee(id uint, storeID uint, updateModels *types.UpdateModels) error {
	if updateModels.StoreEmployee == nil {
		return employeesTypes.ErrNothingToUpdate
	}
	return r.db.Model(&data.StoreEmployee{}).
		Where("id = ? AND store_id = ?", id, storeID).
		Updates(updateModels.StoreEmployee).Error
}
