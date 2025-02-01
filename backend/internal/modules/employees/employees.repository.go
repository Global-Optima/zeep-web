package employees

import (
	"errors"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(employee *data.Employee) (uint, error)

	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	UpdateEmployee(id uint, updateModels *types.UpdateEmployeeModels) error
	UpdateEmployeeWithAssociations(tx *gorm.DB, id uint, updateModels *types.UpdateEmployeeModels) error
	ReassignEmployeeType(employeeID uint, newType data.EmployeeType, workplaceID uint) error

	DeleteTypedEmployeeById(employeeID, workplaceID uint, employeeType data.EmployeeType) error
	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)

	GetEmployeeWorkdayByEmployeeAndDay(employeeID uint, day data.Weekday) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdayByID(workdayID uint) (*data.EmployeeWorkday, error)
	GetEmployeeWorkdaysByEmployeeID(employeeID uint) ([]data.EmployeeWorkday, error)

	GetAllStoreEmployees(storeID uint) ([]data.Employee, error)
	GetAllRegionEmployees(regionID uint) ([]data.Employee, error)
	GetAllFranchiseeEmployees(franchiseeID uint) ([]data.Employee, error)
	GetAllWarehouseEmployees(warehouseID uint) ([]data.Employee, error)
	GetAllAdminEmployees() ([]data.Employee, error)
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
		Preload("StoreEmployee").
		Preload("WarehouseEmployee").
		Preload("RegionEmployee").
		Preload("FranchiseeEmployee").
		Preload("AdminEmployee").
		First(&employee, employeeID).Error
	return &employee, err
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

	if updateModels == nil {
		return types.ErrNothingToUpdate
	}

	if updateModels.Employee != nil {
		err := tx.Model(&data.Employee{}).Where("id = ?", id).Updates(updateModels.Employee).Error
		if err != nil {
			return err
		}
	}

	if updateModels.Workdays == nil {
		err := tx.Where("employees.id = ?", id).Delete(&data.EmployeeWorkday{}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&data.EmployeeWorkday{}).Create(&updateModels.Workdays).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *employeeRepository) ReassignEmployeeType(employeeID uint, newType data.EmployeeType, workplaceID uint) error {
	employee, err := r.GetEmployeeByID(employeeID)
	if err != nil {
		return err
	}

	typeMappings := map[data.EmployeeType]struct {
		deleteModel interface{}
		createModel interface{}
	}{
		data.StoreEmployeeType: {
			deleteModel: &data.StoreEmployee{},
			createModel: &data.StoreEmployee{EmployeeID: employeeID, StoreID: workplaceID},
		},
		data.WarehouseEmployeeType: {
			deleteModel: &data.WarehouseEmployee{},
			createModel: &data.WarehouseEmployee{EmployeeID: employeeID, WarehouseID: workplaceID},
		},
		data.RegionEmployeeType: {
			deleteModel: &data.RegionEmployee{},
			createModel: &data.RegionEmployee{EmployeeID: employeeID, RegionID: workplaceID},
		},
		data.FranchiseeEmployeeType: {
			deleteModel: &data.FranchiseeEmployee{},
			createModel: &data.FranchiseeEmployee{EmployeeID: employeeID, FranchiseeID: workplaceID},
		},
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if deleteModel, ok := typeMappings[employee.GetType()]; ok {
			if err := tx.Where("employee_id = ?", employeeID).Delete(deleteModel.deleteModel).Error; err != nil {
				return err
			}
		} else {
			return types.ErrUnsupportedEmployeeType
		}

		if err := tx.Model(&data.Employee{}).Where("id = ?", employeeID).Update("type", newType).Error; err != nil {
			return err
		}

		if newTypeMapping, ok := typeMappings[newType]; ok {
			if err := tx.Create(newTypeMapping.createModel).Error; err != nil {
				return err
			}
		} else {
			return types.ErrUnsupportedEmployeeType
		}
		return nil
	})
	if err != nil {
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

func (r *employeeRepository) GetAllStoreEmployees(storeID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN store_employees ON store_employees.employee_id = employees.id").
		Where("store_employees.store_id = ?", storeID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllWarehouseEmployees(warehouseID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN warehouse_employees ON warehouse_employees.employee_id = employees.id").
		Where("warehouse_employees.warehouse_id = ?", warehouseID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllAdminEmployees() ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN admin_employees ON admin_employees.employee_id = employees.id").
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllRegionEmployees(regionID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN region_employees ON region_employees.employee_id = employees.id").
		Where("region_employees.region_id = ?", regionID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) GetAllFranchiseeEmployees(franchiseeID uint) ([]data.Employee, error) {
	var employees []data.Employee

	err := r.db.Model(&data.Employee{}).
		Joins("INNER JOIN franchisee_employees ON franchisee_employees.employee_id = employees.id").
		Where("franchisee_employees.franchisee_id = ?", franchiseeID).
		Find(&employees).Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}
