package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type FranchiseeEmployeeRepository interface {
	GetFranchiseeEmployees(franchiseeID uint, filter *employeesTypes.EmployeesFilter) ([]data.FranchiseeEmployee, error)
	GetFranchiseeEmployeeByID(id uint, franchiseeID *uint) (*data.FranchiseeEmployee, error)
	GetAllFranchiseeEmployees(franchiseeID uint) ([]data.FranchiseeEmployee, error)
	UpdateFranchiseeEmployee(id uint, franchiseeID *uint, updateModels *types.UpdateFranchiseeEmployeeModels) error
}

type franchiseeEmployeeRepository struct {
	db           *gorm.DB
	employeeRepo employees.EmployeeRepository
}

func NewFranchiseeEmployeeRepository(db *gorm.DB, employeeRepo employees.EmployeeRepository) FranchiseeEmployeeRepository {
	return &franchiseeEmployeeRepository{
		db:           db,
		employeeRepo: employeeRepo,
	}
}

func (r *franchiseeEmployeeRepository) GetFranchiseeEmployees(franchiseeID uint, filter *employeesTypes.EmployeesFilter) ([]data.FranchiseeEmployee, error) {
	var franchiseeEmployees []data.FranchiseeEmployee
	query := r.db.Model(&data.FranchiseeEmployee{}).
		Where(&data.FranchiseeEmployee{
			FranchiseeID: franchiseeID,
		}).
		Preload("Employee").
		Joins("JOIN employees ON employees.id = franchisee_employees.employee_id")

	if filter.IsActive != nil {
		query = query.Where(&data.FranchiseeEmployee{
			Employee: data.Employee{IsActive: filter.IsActive},
		})
	}

	if filter.Role != nil {
		query = query.Where(&data.FranchiseeEmployee{
			Role: *filter.Role,
		})
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.phone ILIKE ? OR employees.email ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.FranchiseeEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&franchiseeEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve franchisee employees: %w", err)
	}
	return franchiseeEmployees, nil
}

func (r *franchiseeEmployeeRepository) GetFranchiseeEmployeeByID(id uint, franchiseeID *uint) (*data.FranchiseeEmployee, error) {
	var franchiseeEmployee data.FranchiseeEmployee
	query := r.db.Model(&data.FranchiseeEmployee{}).
		Preload("Employee.Workdays").
		Preload("Franchisee").
		Where(&data.FranchiseeEmployee{
			BaseEntity: data.BaseEntity{
				ID: id,
			},
		})

	if franchiseeID != nil {
		query = query.Where(&data.FranchiseeEmployee{
			FranchiseeID: *franchiseeID,
		})
	}

	if err := query.First(&franchiseeEmployee).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve franchisee employee by ID: %w", err)
	}
	return &franchiseeEmployee, nil
}

func (r *franchiseeEmployeeRepository) GetAllFranchiseeEmployees(franchiseeID uint) ([]data.FranchiseeEmployee, error) {
	var franchiseeEmployees []data.FranchiseeEmployee

	err := r.db.Model(&data.FranchiseeEmployee{}).
		Preload("Employee").
		Joins("INNER JOIN employees ON franchisee_employees.employee_id = employees.id").
		Where("franchisee_employees.franchisee_id = ? AND employees.is_active = true", franchiseeID).
		Find(&franchiseeEmployees).Error

	if err != nil {
		return nil, err
	}

	return franchiseeEmployees, nil
}

func (r *franchiseeEmployeeRepository) UpdateFranchiseeEmployee(id uint, franchiseeID *uint, updateModels *types.UpdateFranchiseeEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	existingFranchiseeEmployee, err := r.GetFranchiseeEmployeeByID(id, franchiseeID)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if updateModels.FranchiseeEmployee != nil && !utils.IsEmpty(updateModels.FranchiseeEmployee) {
			err := tx.Model(&data.FranchiseeEmployee{}).
				Where(&data.FranchiseeEmployee{
					BaseEntity: data.BaseEntity{ID: id},
				}).
				Updates(updateModels.FranchiseeEmployee).Error
			if err != nil {
				return err
			}
		}

		if updateModels.UpdateEmployeeModels != nil && !utils.IsEmpty(updateModels.UpdateEmployeeModels) {
			err := r.employeeRepo.UpdateEmployeeWithAssociations(tx, existingFranchiseeEmployee.EmployeeID, updateModels.UpdateEmployeeModels)
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
