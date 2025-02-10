package employees

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegionEmployeeService interface {
	CreateRegionEmployee(regionID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error)
	GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]types.RegionEmployeeDTO, error)
	GetRegionEmployeeByID(id, regionID uint) (*types.RegionEmployeeDetailsDTO, error)
	GetAllRegionEmployees(regionID uint) ([]employeesTypes.EmployeeAccountDTO, error)
	UpdateRegionEmployee(id, regionID uint, input *types.UpdateRegionEmployeeDTO, role data.EmployeeRole) error
}

type regionEmployeeService struct {
	repo         RegionEmployeeRepository
	employeeRepo employees.EmployeeRepository
	logger       *zap.SugaredLogger
}

func NewRegionEmployeeService(repo RegionEmployeeRepository, employeeRepo employees.EmployeeRepository, logger *zap.SugaredLogger) RegionEmployeeService {
	return &regionEmployeeService{
		repo:         repo,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *regionEmployeeService) CreateRegionEmployee(regionID uint, input *employeesTypes.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToRegionEmployee(regionID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create region manager employee", err)
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
		wrappedErr := utils.WrapError("failed to create region manager employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *regionEmployeeService) GetRegionEmployees(regionID uint, filter *employeesTypes.EmployeesFilter) ([]types.RegionEmployeeDTO, error) {
	regionEmployees, err := s.repo.GetRegionEmployees(regionID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve region managers", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.RegionEmployeeDTO, len(regionEmployees))
	for i, regionEmployee := range regionEmployees {
		dtos[i] = *types.MapToRegionEmployeeDTO(&regionEmployee)
	}
	return dtos, nil
}

func (s *regionEmployeeService) GetRegionEmployeeByID(id, regionID uint) (*types.RegionEmployeeDetailsDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid region manager ID")
	}

	employee, err := s.repo.GetRegionEmployeeByID(id, regionID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve region manager with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToRegionEmployeeDetailsDTO(employee), nil
}

func (s *regionEmployeeService) GetAllRegionEmployees(regionID uint) ([]employeesTypes.EmployeeAccountDTO, error) {
	if regionID == 0 {
		return nil, utils.WrapError("invalid region ID", employeesTypes.ErrValidation)
	}

	regionEmployees, err := s.repo.GetAllRegionEmployees(regionID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all region employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]employeesTypes.EmployeeAccountDTO, len(regionEmployees))
	for i, regionEmployee := range regionEmployees {
		dtos[i] = *employeesTypes.MapToEmployeeAccountDTO(&regionEmployee.Employee)
	}

	return dtos, nil
}

func (s *regionEmployeeService) UpdateRegionEmployee(id, regionID uint, input *types.UpdateRegionEmployeeDTO, role data.EmployeeRole) error {
	updateFields, err := types.RegionEmployeeUpdateFields(input, role)
	if err != nil {
		return err
	}
	return s.repo.UpdateRegionEmployee(id, regionID, updateFields)
}
