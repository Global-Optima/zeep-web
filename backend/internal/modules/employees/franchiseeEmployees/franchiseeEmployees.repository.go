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
	GetFranchiseeEmployeeByID(id, franchiseeID uint) (*data.FranchiseeEmployee, error)
	UpdateFranchiseeEmployee(id uint, franchiseeID uint, updateModels *types.UpdateFranchiseeEmployeeModels) error
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
		Where("franchisee_id = ?", franchiseeID).Preload("Employee")

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

func (r *franchiseeEmployeeRepository) GetFranchiseeEmployeeByID(id, franchiseeID uint) (*data.FranchiseeEmployee, error) {
	var franchiseeEmployee data.FranchiseeEmployee
	err := r.db.Model(&data.FranchiseeEmployee{}).
		Preload("Employee").
		Where("id = ? AND franchisee_id = ?", id, franchiseeID).
		First(&franchiseeEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve franchisee employee by ID: %w", err)
	}
	return &franchiseeEmployee, nil
}

func (r *franchiseeEmployeeRepository) UpdateFranchiseeEmployee(id uint, franchiseeID uint, updateModels *types.UpdateFranchiseeEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {

		if updateModels.FranchiseeEmployee != nil {
			err := tx.Model(&data.StoreEmployee{}).
				Where("id = ? AND franchisee_id = ?", id, franchiseeID).
				Updates(updateModels.FranchiseeEmployee).Error
			if err != nil {
				return err
			}
		}

		if updateModels.UpdateEmployeeModels != nil {
			err := r.employeeRepo.UpdateEmployeeWithAssociations(tx, updateModels.FranchiseeEmployee.EmployeeID, updateModels.UpdateEmployeeModels)
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
