package employees

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StoreEmployeeService interface {
	CreateStoreEmployee(storeID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error)
	GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]types.StoreEmployeeDTO, error)
	GetStoreEmployeeByID(id, storeID uint) (*types.StoreEmployeeDTO, error)
	UpdateStoreEmployee(id, storeID uint, input *types.UpdateStoreEmployeeDTO, role data.EmployeeRole) error
}

type storeEmployeeService struct {
	repo         StoreEmployeeRepository
	employeeRepo employees.EmployeeRepository
	logger       *zap.SugaredLogger
}

func NewStoreEmployeeService(repo StoreEmployeeRepository, employeeRepo employees.EmployeeRepository, logger *zap.SugaredLogger) StoreEmployeeService {
	return &storeEmployeeService{
		repo:         repo,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *storeEmployeeService) CreateStoreEmployee(storeID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToStoreEmployee(storeID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store employee", err)
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
		wrappedErr := utils.WrapError("failed to create store employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *storeEmployeeService) GetStoreEmployees(storeID uint, filter *employeesTypes.EmployeesFilter) ([]types.StoreEmployeeDTO, error) {
	storeEmployees, err := s.repo.GetStoreEmployees(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreEmployeeDTO, len(storeEmployees))
	for i, storeEmployee := range storeEmployees {
		dtos[i] = *types.MapToStoreEmployeeDTO(&storeEmployee)
	}

	return dtos, nil
}

func (s *storeEmployeeService) GetStoreEmployeeByID(id, storeID uint) (*types.StoreEmployeeDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid store employee ID")
	}

	employee, err := s.repo.GetStoreEmployeeByID(id, storeID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve store employee for store employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToStoreEmployeeDTO(employee), nil
}

func (s *storeEmployeeService) UpdateStoreEmployee(id, storeID uint, input *types.UpdateStoreEmployeeDTO, role data.EmployeeRole) error {
	updateFields, err := types.StoreEmployeeUpdateFields(input, role)
	if err != nil {
		return err
	}
	return s.repo.UpdateStoreEmployee(id, storeID, updateFields)
}
