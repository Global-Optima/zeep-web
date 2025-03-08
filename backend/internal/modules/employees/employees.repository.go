package employees

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(employee *data.Employee) (uint, error)

	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	GetEmployees(filter *types.EmployeesFilter) ([]data.Employee, error)
	UpdateEmployee(id uint, updateModels *types.UpdateEmployeeModels) error
	UpdateEmployeeWithAssociations(tx *gorm.DB, id uint, updateModels *types.UpdateEmployeeModels) error
	ReassignEmployeeType(employeeID uint, dto *types.ReassignEmployeeTypeDTO) error

	DeleteTypedEmployeeById(employeeID, workplaceID uint, employeeType data.EmployeeType) error
	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)

	GetEmployeeWorkdayByEmployeeAndDay(employeeID uint, day data.Weekday) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdayByID(workdayID uint) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdaysByEmployeeID(employeeID uint) ([]data.EmployeeWorkday, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(employee *data.Employee) (uint, error) {
	err := r.db.Create(employee).Error
	if err != nil {
		return 0, err
	}
	return employee.ID, nil
}

func (r *employeeRepository) GetEmployeeByID(employeeID uint) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Model(&data.Employee{}).
		Preload("StoreEmployee.Store").
		Preload("WarehouseEmployee.Warehouse").
		Preload("RegionEmployee.Region").
		Preload("FranchiseeEmployee.Franchisee").
		Preload("AdminEmployee").
		Preload("Workdays").
		First(&employee, employeeID).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, moduleErrors.ErrNotFound
	}

	return &employee, err
}

func (r *employeeRepository) GetEmployees(filter *types.EmployeesFilter) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Model(&data.Employee{})

	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	if filter.Role != nil {
		employeeType := data.GetEmployeeTypeByRole(*filter.Role)
		switch employeeType {
		case data.WarehouseEmployeeType:
			query.Joins("JOIN warehouse_employees ON warehouse_employees.employee_id = employees.id").
				Where("warehouse_employees.role = ?", *filter.Role).
				Preload("WarehouseEmployee")
		case data.StoreEmployeeType:
			query.Joins("JOIN store_employees ON store_employees.employee_id = employees.id").
				Where("store_employees.role = ?", *filter.Role).
				Preload("StoreEmployee")
		case data.RegionEmployeeType:
			query.Joins("JOIN region_employees ON region_employees.employee_id = employees.id").
				Where("region_employees.role = ?", *filter.Role).
				Preload("RegionEmployee")
		case data.FranchiseeEmployeeType:
			query.Joins("JOIN franchisee_employees ON franchisee_employees.employee_id = employees.id").
				Where("franchisee_employees.role = ?", *filter.Role).
				Preload("FranchiseeEmployee")
		case data.AdminEmployeeType:
			query.Joins("JOIN store_employees ON store_employees.employee_id = employees.id").
				Where("store_employees.role = ?", *filter.Role).
				Preload("AdminEmployee")
		default:
			return nil, fmt.Errorf("%w: %s", types.ErrUnsupportedEmployeeType, employeeType)
		}
	} else {
		query.Preload("StoreEmployee").
			Preload("WarehouseEmployee").
			Preload("RegionEmployee").
			Preload("FranchiseeEmployee").
			Preload("AdminEmployee")
	}

	if filter.IsActive != nil {
		query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query.Where("(first_name || ' ' || last_name) ILIKE ?", "%"+searchTerm+"%")
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Employee{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&employees).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, moduleErrors.ErrNotFound
	}

	return employees, err
}

func (r *employeeRepository) GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Model(&data.Employee{}).
		Preload("StoreEmployee").
		Preload("WarehouseEmployee").
		Preload("FranchiseeEmployee").
		Preload("RegionEmployee").
		Preload("AdminEmployee").
		Where("email = ? OR phone = ?", email, phone).
		First(&employee).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) UpdateEmployee(id uint, updateModels *types.UpdateEmployeeModels) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return r.UpdateEmployeeWithAssociations(tx, id, updateModels)
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *employeeRepository) UpdateEmployeeWithAssociations(tx *gorm.DB, id uint, updateModels *types.UpdateEmployeeModels) error {
	if updateModels.Employee != nil && !utils.IsEmpty(updateModels.Employee) {
		err := tx.Model(&data.Employee{}).Where("id = ?", id).Updates(updateModels.Employee).Error
		if err != nil {
			return err
		}
	}

	if updateModels.Workdays != nil && !utils.IsEmpty(updateModels.Workdays) {
		err := tx.Where("employee_id = ?", id).Unscoped().Delete(&data.EmployeeWorkday{}).Error
		if err != nil {
			return err
		}

		for i := 0; i < len(updateModels.Workdays); i++ {
			updateModels.Workdays[i].EmployeeID = id
		}

		err = tx.Model(&data.EmployeeWorkday{}).Create(&updateModels.Workdays).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *employeeRepository) ReassignEmployeeType(employeeID uint, dto *types.ReassignEmployeeTypeDTO) error {
	employee, err := r.GetEmployeeByID(employeeID)
	if err != nil {
		return err
	}

	typeMappings := map[data.EmployeeType]struct {
		sampleModel interface{}
		createModel interface{}
	}{
		data.StoreEmployeeType: {
			sampleModel: &data.StoreEmployee{},
			createModel: &data.StoreEmployee{EmployeeID: employeeID, Role: dto.Role, StoreID: dto.WorkplaceID},
		},
		data.WarehouseEmployeeType: {
			sampleModel: &data.WarehouseEmployee{},
			createModel: &data.WarehouseEmployee{EmployeeID: employeeID, Role: dto.Role, WarehouseID: dto.WorkplaceID},
		},
		data.RegionEmployeeType: {
			sampleModel: &data.RegionEmployee{},
			createModel: &data.RegionEmployee{EmployeeID: employeeID, Role: dto.Role, RegionID: dto.WorkplaceID},
		},
		data.FranchiseeEmployeeType: {
			sampleModel: &data.FranchiseeEmployee{},
			createModel: &data.FranchiseeEmployee{EmployeeID: employeeID, Role: dto.Role, FranchiseeID: dto.WorkplaceID},
		},
	}

	currentEmployeeType := employee.GetType()
	models, ok := typeMappings[currentEmployeeType]
	if !ok {
		return types.ErrUnsupportedEmployeeType
	}

	if currentEmployeeType != dto.EmployeeType {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			{
				if err := tx.Where("employee_id = ?", employeeID).Delete(models.sampleModel).Error; err != nil {
					return err
				}
			}
			if newTypeMapping, ok := typeMappings[dto.EmployeeType]; ok {
				if err := tx.Create(newTypeMapping.createModel).Error; err != nil {
					return err
				}
			} else {
				return types.ErrUnsupportedEmployeeType
			}
			return nil
		})
		return err
	}

	if err := r.db.Where("employee_id = ?", employeeID).Updates(models.createModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *employeeRepository) DeleteTypedEmployeeById(employeeID, workplaceID uint, employeeType data.EmployeeType) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND type = ?", employeeID, employeeType).Delete(&data.Employee{}).Error; err != nil {
			return err
		}

		switch employeeType {
		case data.StoreEmployeeType:
			if err := tx.Where("id = ? AND store_id", employeeID, workplaceID).Delete(&data.StoreEmployee{}).Error; err != nil {
				return err
			}
		case data.WarehouseEmployeeType:
			if err := tx.Where("id = ? AND warehouse_id", employeeID, workplaceID).Delete(&data.WarehouseEmployee{}).Error; err != nil {
				return err
			}
		case data.FranchiseeEmployeeType:
			if err := tx.Where("id = ? AND franchisee_id", employeeID, workplaceID).Delete(&data.FranchiseeEmployee{}).Error; err != nil {
				return err
			}
		case data.RegionEmployeeType:
			if err := tx.Where("id = ? AND region_id", employeeID, workplaceID).Delete(&data.RegionEmployee{}).Error; err != nil {
				return err
			}
		default:
			return types.ErrUnsupportedEmployeeType
		}

		if err := tx.Where("employee_id = ?", employeeID).Delete(&data.EmployeeWorkday{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *employeeRepository) GetEmployeeWorkdayByEmployeeAndDay(employeeID uint, day data.Weekday) (*data.EmployeeWorkday, error) {
	var workday data.EmployeeWorkday

	if employeeID == 0 {
		return nil, fmt.Errorf("employeeID cannot be zero")
	}

	err := r.db.Where("employee_id = ? AND day = ?", employeeID, day).First(&workday).Error
	if err != nil {
		return nil, err
	}
	return &workday, err
}

func (r *employeeRepository) GetEmployeeWorkdayByID(workdayID uint) (*data.EmployeeWorkday, error) {
	var workday data.EmployeeWorkday
	err := r.db.
		Preload("Employee").
		First(&workday, workdayID).Error
	if err != nil {
		return nil, err
	}
	return &workday, nil
}

func (r *employeeRepository) GetEmployeeWorkdaysByEmployeeID(employeeID uint) ([]data.EmployeeWorkday, error) {
	var workdays []data.EmployeeWorkday
	err := r.db.Where("employee_id = ?", employeeID).Find(&workdays).Error
	if err != nil {
		return nil, err
	}
	return workdays, nil
}
