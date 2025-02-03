package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type RegionEmployeeRepository interface {
	GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]data.RegionEmployee, error)
	GetRegionEmployeeByID(id, regionID uint) (*data.RegionEmployee, error)
	UpdateRegionEmployee(id uint, regionID uint, updateModels *types.UpdateRegionEmployeeModels) error
}

type regionEmployeeRepository struct {
	db           *gorm.DB
	employeeRepo employees.EmployeeRepository
}

func NewRegionEmployeeRepository(db *gorm.DB, employeeRepo employees.EmployeeRepository) RegionEmployeeRepository {
	return &regionEmployeeRepository{
		db:           db,
		employeeRepo: employeeRepo,
	}
}

func (r *regionEmployeeRepository) GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]data.RegionEmployee, error) {
	var regionEmployees []data.RegionEmployee
	query := r.db.Model(&data.RegionEmployee{}).
		Where("region_id = ?", regionID).Preload("Employee")

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

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.RegionEmployee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&regionEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region managers: %w", err)
	}
	return regionEmployees, nil
}

func (r *regionEmployeeRepository) GetRegionEmployeeByID(id, regionID uint) (*data.RegionEmployee, error) {
	var regionEmployee data.RegionEmployee
	err := r.db.Model(&data.RegionEmployee{}).
		Preload("Employee.Workdays").
		Preload("Region").
		Where("id = ? AND region_id = ?", id, regionID).
		First(&regionEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region manager by ID: %w", err)
	}
	return &regionEmployee, nil
}

// TODO add employee update to all submodules
func (r *regionEmployeeRepository) UpdateRegionEmployee(id uint, regionID uint, updateModels *types.UpdateRegionEmployeeModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existingRegionEmployee data.RegionEmployee
		r.db.Model(&data.RegionEmployee{}).
			Where("id = ? AND region_id = ?", id, regionID).
			First(&existingRegionEmployee)

		if updateModels.RegionEmployee != nil {
			err := tx.Model(&data.RegionEmployee{}).
				Where("id = ? AND region_id = ?", id, regionID).
				Updates(updateModels.RegionEmployee).Error
			if err != nil {
				return err
			}
		}

		if updateModels.UpdateEmployeeModels != nil {
			err := r.employeeRepo.UpdateEmployeeWithAssociations(tx, existingRegionEmployee.EmployeeID, updateModels.UpdateEmployeeModels)
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
