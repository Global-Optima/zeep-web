package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type RegionEmployeeRepository interface {
	GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]data.RegionEmployee, error)
	GetRegionEmployeeByID(id, regionID uint) (*data.RegionEmployee, error)
	UpdateRegionEmployee(id uint, regionID uint, updateModels *types.UpdateModels) error
}

type regionEmployeeRepository struct {
	db *gorm.DB
}

func NewRegionEmployeeRepository(db *gorm.DB) RegionEmployeeRepository {
	return &regionEmployeeRepository{db: db}
}

func (r *regionEmployeeRepository) GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]data.RegionEmployee, error) {
	var employees []data.RegionEmployee
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

	err = query.Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region managers: %w", err)
	}
	return employees, nil
}

func (r *regionEmployeeRepository) GetRegionEmployeeByID(id, regionID uint) (*data.RegionEmployee, error) {
	var regionEmployee data.RegionEmployee
	err := r.db.Model(&data.RegionEmployee{}).
		Preload("Employee").
		Where("id = ? AND region_id = ?", id, regionID).
		First(&regionEmployee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region manager by ID: %w", err)
	}
	return &regionEmployee, nil
}

// TODO add employee update to all submodules
func (r *regionEmployeeRepository) UpdateRegionEmployee(id uint, regionID uint, updateModels *types.UpdateModels) error {
	if updateModels == nil {
		return employeesTypes.ErrNothingToUpdate
	}
	return r.db.Model(&data.RegionEmployee{}).
		Where("id = ? AND region_id = ?", id, regionID).
		Updates(updateModels.RegionEmployee).Error
}
