package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StoreEmployeeRepository interface {
	GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]data.StoreEmployee, error)
	GetStoreEmployeeByID(id, storeID uint) (*data.StoreEmployee, error)
	GetAllStoreEmployees(storeID uint) ([]data.StoreEmployee, error)
	UpdateStoreEmployee(id uint, storeID uint, updateModels *types.UpdateStoreEmployeeModels) error
}

type storeEmployeeRepository struct {
	db           *gorm.DB
	employeeRepo employees.EmployeeRepository
}

func NewStoreEmployeeRepository(db *gorm.DB, employeeRepo employees.EmployeeRepository) StoreEmployeeRepository {
	return &storeEmployeeRepository{
		db:           db,
		employeeRepo: employeeRepo,
	}
}

func (r *storeEmployeeRepository) GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]data.StoreEmployee, error) {
	var storeEmployees []data.StoreEmployee
	query := r.db.Model(&data.StoreEmployee{}).
		Where("store_id = ?", storeID).
		Preload("Employee").
		Joins("JOIN employees ON employees.id = store_employees.employee_id")

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Role != nil {
		query = query.Where("role = ?", *filter.Role)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&storeEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employees: %w", err)
	}
	return storeEmployees, nil
}

func (r *storeEmployeeRepository) GetStoreEmployeeByID(id, storeID uint) (*data.StoreEmployee, error) {
	var storeEmployee data.StoreEmployee
	err := r.db.Model(&data.StoreEmployee{}).
		Preload("Employee.Workdays").
		Preload("Store.FacilityAddress").
		Where("id = ? AND store_id = ?", id, storeID).
		First(&storeEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store employee by ID: %w", err)
	}
	return &storeEmployee, nil
}

func (r *storeEmployeeRepository) UpdateStoreEmployee(id uint, storeID uint, updateModels *types.UpdateStoreEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existingStoreEmployee data.StoreEmployee
		r.db.Model(&data.StoreEmployee{}).
			Where("id = ? AND store_id = ?", id, storeID).
			First(&existingStoreEmployee)

		if updateModels.StoreEmployee != nil {
			err := tx.Model(&data.StoreEmployee{}).
				Where("id = ? AND store_id = ?", id, storeID).
				Updates(updateModels.StoreEmployee).Error
			if err != nil {
				return err
			}
		}

		if updateModels.UpdateEmployeeModels != nil {
			err := r.employeeRepo.UpdateEmployeeWithAssociations(tx, existingStoreEmployee.EmployeeID, updateModels.UpdateEmployeeModels)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *storeEmployeeRepository) GetAllStoreEmployees(storeID uint) ([]data.StoreEmployee, error) {
	var storeEmployees []data.StoreEmployee

	err := r.db.Model(&data.StoreEmployee{}).
		Joins("INNER JOIN employees ON store_employees.employee_id = employees.id").
		Where("store_id = ?", storeID).
		Preload("Employee").
		Find(&storeEmployees).Error

	if err != nil {
		return nil, err
	}

	return storeEmployees, nil
}
