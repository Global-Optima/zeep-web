package employees

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FranchiseeEmployeeService interface {
	CreateFranchiseeEmployee(franchiseeID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error)
	GetFranchiseeEmployees(franchiseeID uint, filter *employeesTypes.EmployeesFilter) ([]types.FranchiseeEmployeeDTO, error)
	GetFranchiseeEmployeeByID(id, franchiseeID uint) (*types.FranchiseeEmployeeDetailsDTO, error)
	UpdateFranchiseeEmployee(id, franchiseeID uint, input *types.UpdateFranchiseeEmployeeDTO, role data.EmployeeRole) error
}

type franchiseeEmployeeService struct {
	repo         FranchiseeEmployeeRepository
	employeeRepo employees.EmployeeRepository
	logger       *zap.SugaredLogger
}

func NewFranchiseeEmployeeService(repo FranchiseeEmployeeRepository, employeeRepo employees.EmployeeRepository, logger *zap.SugaredLogger) FranchiseeEmployeeService {
	return &franchiseeEmployeeService{
		repo:         repo,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *franchiseeEmployeeService) CreateFranchiseeEmployee(franchiseeID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToFranchiseeEmployee(franchiseeID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create franchisee employee", err)
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
		wrappedErr := utils.WrapError("failed to create franchisee employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *franchiseeEmployeeService) GetFranchiseeEmployees(franchiseeID uint, filter *employeesTypes.EmployeesFilter) ([]types.FranchiseeEmployeeDTO, error) {
	franchiseeEmployees, err := s.repo.GetFranchiseeEmployees(franchiseeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve franchisee employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.FranchiseeEmployeeDTO, len(franchiseeEmployees))
	for i, franchiseeEmployee := range franchiseeEmployees {
		dtos[i] = *types.MapToFranchiseeEmployeeDTO(&franchiseeEmployee)
	}
	return dtos, nil
}

func (s *franchiseeEmployeeService) GetFranchiseeEmployeeByID(id, franchiseeID uint) (*types.FranchiseeEmployeeDetailsDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid franchise employee ID")
	}

	employee, err := s.repo.GetFranchiseeEmployeeByID(id, franchiseeID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve franchisee employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToFranchiseeEmployeeDetailsDTO(employee), nil
}

func (s *franchiseeEmployeeService) UpdateFranchiseeEmployee(id, franchiseeID uint, input *types.UpdateFranchiseeEmployeeDTO, role data.EmployeeRole) error {
	updateFields, err := types.FranchiseeEmployeeUpdateFields(input, role)
	if err != nil {
		return err
	}
	return s.repo.UpdateFranchiseeEmployee(id, franchiseeID, updateFields)
}
