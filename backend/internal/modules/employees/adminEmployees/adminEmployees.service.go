package employees

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/adminEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminEmployeeService interface {
	CreateAdminEmployee(dto *employeesTypes.CreateEmployeeDTO) (uint, error)
	GetAdminEmployees(filter *employeesTypes.EmployeesFilter) ([]types.AdminEmployeeDTO, error)
	GetAdminEmployeeByID(id uint) (*types.AdminEmployeeDTO, error)
}

type adminEmployeeService struct {
	repo         AdminEmployeeRepository
	employeeRepo employees.EmployeeRepository
	logger       *zap.SugaredLogger
}

func NewAdminEmployeeService(repo AdminEmployeeRepository, employeeRepo employees.EmployeeRepository, logger *zap.SugaredLogger) AdminEmployeeService {
	return &adminEmployeeService{
		repo:         repo,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *adminEmployeeService) CreateAdminEmployee(input *employeesTypes.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToAdminEmployee(input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create admin employee", err)
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
		wrappedErr := utils.WrapError("failed to create admin employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *adminEmployeeService) GetAdminEmployees(filter *employeesTypes.EmployeesFilter) ([]types.AdminEmployeeDTO, error) {
	adminEmployees, err := s.repo.GetAdminEmployees(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve admin employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.AdminEmployeeDTO, len(adminEmployees))
	for i, employee := range adminEmployees {
		dtos[i] = *types.MapToAdminEmployeeDTO(&employee)
	}
	return dtos, nil
}

func (s *adminEmployeeService) GetAdminEmployeeByID(id uint) (*types.AdminEmployeeDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid admin employee ID")
	}

	employee, err := s.repo.GetAdminEmployeeByID(id)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve admin employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToAdminEmployeeDTO(employee), nil
}
