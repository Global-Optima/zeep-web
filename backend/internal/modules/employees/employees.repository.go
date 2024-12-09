package employees

import (
	"errors"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(employee *data.Employee) error
	GetEmployees(query types.GetEmployeesQuery) ([]data.Employee, error)
	GetStoreEmployees(storeID uint, role *string, limit, offset int) ([]data.Employee, error)
	GetWarehouseEmployees(warehouseID uint, role *string, limit, offset int) ([]data.Employee, error)
	GetEmployeeByID(employeeID uint) (*data.Employee, error)
	UpdateEmployee(employee *data.Employee) error
	PartialUpdateEmployee(employeeID uint, fields map[string]interface{}) error
	DeleteEmployee(employeeID uint) error

	GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error)
	GetAllRoles() ([]data.EmployeeRole, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(employee *data.Employee) error {
	return r.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(employee).Error; err != nil {
				return err
			}

			if employee.Type == data.StoreEmployeeType {
				storeEmployee := &data.StoreEmployee{
					EmployeeID:  employee.ID,
					StoreID:     employee.StoreEmployee.StoreID,
					IsFranchise: employee.StoreEmployee.IsFranchise,
				}

				if err := tx.Create(storeEmployee).Error; err != nil {
					return err
				}
			} else if employee.Type == data.WarehouseEmployeeType {
				warehouseEmployee := &data.WarehouseEmployee{
					EmployeeID:  employee.ID,
					WarehouseID: employee.WarehouseEmployee.WarehouseID,
				}
				if err := tx.Create(warehouseEmployee).Error; err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func (r *employeeRepository) GetEmployees(query types.GetEmployeesQuery) ([]data.Employee, error) {
	var employees []data.Employee
	dbQuery := r.db.Where("is_active = TRUE")

	if query.Type != nil {
		dbQuery = dbQuery.Where("LOWER(type) = ?", strings.ToLower(*query.Type))
	}
	if query.StoreID != nil {
		dbQuery = dbQuery.Joins("LEFT JOIN store_employees ON employees.id = store_employees.employee_id").
			Where("store_employees.store_id = ?", *query.StoreID)
	}
	if query.WarehouseID != nil {
		dbQuery = dbQuery.Joins("LEFT JOIN warehouse_employees ON employees.id = warehouse_employees.employee_id").
			Where("warehouse_employees.warehouse_id = ?", *query.WarehouseID)
	}
	if query.Role != nil {
		dbQuery = dbQuery.Where("LOWER(role) = ?", strings.ToLower(*query.Role))
	}

	if query.Limit > 0 {
		dbQuery = dbQuery.Limit(query.Limit)
	}
	if query.Offset >= 0 {
		dbQuery = dbQuery.Offset(query.Offset)
	}

	dbQuery = dbQuery.Preload("StoreEmployee").Preload("WarehouseEmployee")

	err := dbQuery.Find(&employees).Error
	if err != nil {
		return nil, err
	}

	return employees, err
}

func (r *employeeRepository) GetStoreEmployees(storeID uint, role *string, limit, offset int) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Preload("StoreEmployee").Where("is_active = TRUE").
		Joins("JOIN store_employees ON employees.id = store_employees.employee_id").
		Where("store_employees.store_id = ?", storeID)

	if role != nil {
		query = query.Where("employees.role = ?", *role)
	}

	err := query.Limit(limit).Offset(offset).Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) GetWarehouseEmployees(warehouseID uint, role *string, limit, offset int) ([]data.Employee, error) {
	var employees []data.Employee
	query := r.db.Preload("WarehouseEmployee").Where("is_active = TRUE").
		Joins("JOIN warehouse_employees ON employees.id = warehouse_employees.employee_id").
		Where("warehouse_employees.warehouse_id = ?", warehouseID)

	if role != nil {
		query = query.Where("employees.role = ?", *role)
		query = query.Where("employees.role = ?", *role)
	}

	err := query.Limit(limit).Offset(offset).Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) GetEmployeeByID(employeeID uint) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Preload("StoreEmployee").Preload("WarehouseEmployee").First(&employee, employeeID).Error
	return &employee, err
}

func (r *employeeRepository) GetEmployeeByEmailOrPhone(email string, phone string) (*data.Employee, error) {
	var employee data.Employee
	err := r.db.Preload("StoreEmployee").Preload("WarehouseEmployee").
		Where("email = ? OR phone = ?", email, phone).
		First(&employee).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &employee, err
}

func (r *employeeRepository) UpdateEmployee(employee *data.Employee) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(employee).Error; err != nil {
			return err
		}

		if employee.Type == data.StoreEmployeeType && employee.StoreEmployee != nil {
			if err := tx.Save(employee.StoreEmployee).Error; err != nil {
				return err
			}
		} else if employee.Type == data.WarehouseEmployeeType && employee.WarehouseEmployee != nil {
			if err := tx.Save(employee.WarehouseEmployee).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *employeeRepository) PartialUpdateEmployee(employeeID uint, fields map[string]interface{}) error {
	if storeEmployee, ok := fields["store_employee"].(map[string]interface{}); ok {
		if err := r.db.Model(&data.StoreEmployee{}).Where("employee_id = ?", employeeID).Updates(storeEmployee).Error; err != nil {
			return err
		}
		delete(fields, "store_employee")
	}

	if warehouseEmployee, ok := fields["warehouse_employee"].(map[string]interface{}); ok {
		if err := r.db.Model(&data.WarehouseEmployee{}).Where("employee_id = ?", employeeID).Updates(warehouseEmployee).Error; err != nil {
			return err
		}
		delete(fields, "warehouse_employee")
	}

	return r.db.Model(&data.Employee{}).Where("id = ?", employeeID).Updates(fields).Error
}

func (r *employeeRepository) DeleteEmployee(employeeID uint) error {
	return r.db.Model(&data.Employee{}).Where("id = ?", employeeID).Update("is_active", false).Error
}

func (r *employeeRepository) GetAllRoles() ([]data.EmployeeRole, error) {
	roles := []data.EmployeeRole{
		data.RoleAdmin,
		data.RoleManager,
		data.RoleBarista,
	}
	return roles, nil
}
