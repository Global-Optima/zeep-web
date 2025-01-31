package employees

import (
	"fmt"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type FranchiseeEmployeeRepository interface {
	GetFranchiseeEmployees(franchiseeID uint, filter *employeesTypes.EmployeesFilter) ([]data.FranchiseeEmployee, error)
	GetFranchiseeEmployeeByID(id, franchiseeID uint) (*data.FranchiseeEmployee, error)
	UpdateFranchiseeEmployee(id uint, franchiseeID uint, updateModels *types.UpdateModels) error
}

type franchiseeEmployeeRepository struct {
	db *gorm.DB
}

func NewFranchiseeEmployeeRepository(db *gorm.DB) FranchiseeEmployeeRepository {
	return &franchiseeEmployeeRepository{db: db}
}

func (r *franchiseeEmployeeRepository) GetFranchiseeEmployees(franchiseeID uint, filter *employeesTypes.EmployeesFilter) ([]data.FranchiseeEmployee, error) {
	var employees []data.FranchiseeEmployee
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

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve franchisee employees: %w", err)
	}
	return employees, nil
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

func (r *franchiseeEmployeeRepository) UpdateFranchiseeEmployee(id uint, franchiseeID uint, updateModels *types.UpdateModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}
	return r.db.Model(&data.FranchiseeEmployee{}).
		Where("id = ? AND franchisee_id = ?", id, franchiseeID).
		Updates(updateModels.FranchiseeEmployee).Error
}
