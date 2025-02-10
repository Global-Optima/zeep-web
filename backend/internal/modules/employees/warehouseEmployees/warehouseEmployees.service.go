package employees

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/warehouseEmployees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WarehouseEmployeeService interface {
	CreateWarehouseEmployee(warehouseID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error)
	GetWarehouseEmployees(warehouseID uint, filter *employeesTypes.EmployeesFilter) ([]types.WarehouseEmployeeDTO, error)
	GetWarehouseEmployeeByID(id uint, filter *contexts.WarehouseContextFilter) (*types.WarehouseEmployeeDetailsDTO, error)
	GetAllWarehouseEmployees(warehouseID uint) ([]employeesTypes.EmployeeAccountDTO, error)
	UpdateWarehouseEmployee(id uint, filter *contexts.WarehouseContextFilter, input *types.UpdateWarehouseEmployeeDTO, role data.EmployeeRole) error
}

type warehouseEmployeeService struct {
	repo         WarehouseEmployeeRepository
	employeeRepo employees.EmployeeRepository
	logger       *zap.SugaredLogger
}

func NewWarehouseEmployeeService(repo WarehouseEmployeeRepository, employeeRepo employees.EmployeeRepository, logger *zap.SugaredLogger) WarehouseEmployeeService {
	return &warehouseEmployeeService{
		repo:         repo,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *warehouseEmployeeService) CreateWarehouseEmployee(warehouseID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToWarehouseEmployee(warehouseID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create warehouse employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	existingEmployee, err := s.employeeRepo.GetEmployeeByEmailOrPhone(employee.Email, employee.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking employee uniqueness", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	if existingEmployee != nil {
		return 0, employeesTypes.ErrEmployeeAlreadyExists
	}

	id, err := s.employeeRepo.CreateEmployee(employee)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create warehouse employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *warehouseEmployeeService) GetWarehouseEmployees(warehouseID uint, filter *employeesTypes.EmployeesFilter) ([]types.WarehouseEmployeeDTO, error) {
	warehouseEmployees, err := s.repo.GetWarehouseEmployees(warehouseID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.WarehouseEmployeeDTO, len(warehouseEmployees))
	for i, warehouseEmployee := range warehouseEmployees {
		warehouseEmployee.Employee.WarehouseEmployee = &warehouseEmployee
		dtos[i] = *types.MapToWarehouseEmployeeDTO(&warehouseEmployee)
	}

	return dtos, nil
}

func (s *warehouseEmployeeService) GetWarehouseEmployeeByID(id uint, filter *contexts.WarehouseContextFilter) (*types.WarehouseEmployeeDetailsDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid warehouse employee ID")
	}

	employee, err := s.repo.GetWarehouseEmployeeByID(id, filter)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve warehouse employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToWarehouseEmployeeDetailsDTO(employee), nil
}

func (s *warehouseEmployeeService) GetAllWarehouseEmployees(warehouseID uint) ([]employeesTypes.EmployeeAccountDTO, error) {
	if warehouseID == 0 {
		return nil, utils.WrapError("invalid warehouse ID", employeesTypes.ErrValidation)
	}

	warehouseEmployees, err := s.repo.GetAllWarehouseEmployees(warehouseID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all warehouse employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]employeesTypes.EmployeeAccountDTO, len(warehouseEmployees))
	for i, warehouseEmployee := range warehouseEmployees {
		dtos[i] = *employeesTypes.MapToEmployeeAccountDTO(&warehouseEmployee.Employee)
	}

	return dtos, nil
}

func (s *warehouseEmployeeService) UpdateWarehouseEmployee(id uint, filter *contexts.WarehouseContextFilter, input *types.UpdateWarehouseEmployeeDTO, role data.EmployeeRole) error {
	updateFields, err := types.WarehouseEmployeeUpdateFields(input, role)
	if err != nil {
		return err
	}
	return s.repo.UpdateWarehouseEmployee(id, filter, updateFields)
}
