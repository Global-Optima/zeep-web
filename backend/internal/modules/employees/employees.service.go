package employees

import (
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EmployeeService interface {
	CreateEmployee(employee *data.Employee) (uint, error)

	GetEmployeeByID(id uint) (*types.EmployeeDetailsDTO, error)
	GetEmployees(filter *types.EmployeesFilter) ([]types.EmployeeDTO, error)
	UpdateEmployeeInfo(employeeID uint, dto *types.UpdateEmployeeDTO) error
	ReassignEmployeeType(employeeID uint, dto *types.ReassignEmployeeTypeDTO) error
	DeleteTypedEmployee(employeeID, workplaceID uint, employeeType data.EmployeeType) error
	UpdatePassword(employeeID uint, input *types.UpdatePasswordDTO) error

	GetAllRoles() ([]types.EmployeeTypeRoles, error)
	GetEmployeeWorkday(workdayID uint) (*types.EmployeeWorkdayDTO, error)
	GetEmployeeWorkdays(employeeID uint) ([]types.EmployeeWorkdayDTO, error)

	GetAllWarehouseEmployees(warehouseID uint) ([]types.EmployeeAccountDTO, error)
	GetAllStoreEmployees(storeID uint) ([]types.EmployeeAccountDTO, error)
	GetAllRegionEmployees(regionID uint) ([]types.EmployeeAccountDTO, error)
	GetAllFranchiseeEmployees(franchiseeID uint) ([]types.EmployeeAccountDTO, error)
	GetAllAdminEmployees() ([]types.EmployeeAccountDTO, error)
}

type employeeService struct {
	repo   EmployeeRepository
	logger *zap.SugaredLogger
}

func NewEmployeeService(repo EmployeeRepository, logger *zap.SugaredLogger) EmployeeService {
	return &employeeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *employeeService) CreateEmployee(employee *data.Employee) (uint, error) {
	existingEmployee, err := s.repo.GetEmployeeByEmailOrPhone(employee.Email, employee.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError("error checking employee uniqueness", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	if existingEmployee != nil {
		return 0, types.ErrEmployeeAlreadyExists
	}

	id, err := s.repo.CreateEmployee(employee)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *employeeService) GetEmployeeByID(id uint) (*types.EmployeeDetailsDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid store employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(id)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToEmployeeDetailsDTO(employee), nil
}

func (s *employeeService) GetEmployees(filter *types.EmployeesFilter) ([]types.EmployeeDTO, error) {
	employees, err := s.repo.GetEmployees(filter)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve employees")
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employees == nil {
		return nil, errors.New("employee not found")
	}

	dtos := make([]types.EmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToEmployeeDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) UpdateEmployeeInfo(employeeID uint, dto *types.UpdateEmployeeDTO) error {
	if employeeID == 0 {
		return fmt.Errorf("%w: invalid employee ID: %d", types.ErrValidation, employeeID)
	}

	updateModels, err := types.PrepareUpdateFields(dto)
	if err != nil {
		wrappedErr := fmt.Errorf("%w: %v", types.ErrValidation, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	err = s.repo.UpdateEmployee(employeeID, updateModels)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *employeeService) ReassignEmployeeType(employeeID uint, dto *types.ReassignEmployeeTypeDTO) error {
	if employeeID == 0 {
		return fmt.Errorf("%w: invalid employee ID: %d", types.ErrValidation, employeeID)
	}

	err := s.repo.ReassignEmployeeType(employeeID, dto.EmployeeType, dto.WorkplaceID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to reassign employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *employeeService) DeleteTypedEmployee(employeeID, workplaceID uint, employeeType data.EmployeeType) error {
	if employeeID == 0 {
		return errors.New("invalid employee ID")
	}

	if err := s.repo.DeleteTypedEmployeeById(employeeID, workplaceID, employeeType); err != nil {
		wrappedErr := fmt.Errorf("failed to delete %s employee with ID = %d: %w", employeeType, employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *employeeService) UpdatePassword(employeeID uint, input *types.UpdatePasswordDTO) error {
	if employeeID == 0 {
		return errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(employeeID)
	if err != nil {
		return fmt.Errorf("failed to retrieve employee: %v", err)
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	if err := utils.ComparePassword(employee.HashedPassword, input.OldPassword); err != nil {
		return errors.New("incorrect old password")
	}

	if err := utils.IsValidPassword(input.NewPassword); err != nil {
		return fmt.Errorf("password validation failed: %v", err)
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	updateModels := types.UpdateEmployeeModels{
		Employee: employee,
	}

	employee.HashedPassword = hashedPassword
	if err := s.repo.UpdateEmployee(employeeID, &updateModels); err != nil {
		wrappedErr := fmt.Errorf("failed to update password for employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *employeeService) GetAllRoles() ([]types.EmployeeTypeRoles, error) {
	var employeeTypeRoles []types.EmployeeTypeRoles

	for employeeType, roles := range data.EmployeeTypeRoleMap {
		employeeTypeRoles = append(employeeTypeRoles, types.EmployeeTypeRoles{
			EmployeeType: employeeType,
			Roles:        roles,
		})
	}

	return employeeTypeRoles, nil
}

func (s *employeeService) GetEmployeeWorkday(workdayID uint) (*types.EmployeeWorkdayDTO, error) {
	workday, err := s.repo.GetEmployeeWorkdayByID(workdayID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employee workday", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dto := types.MapToEmployeeWorkdayDTO(workday)
	return dto, nil
}

func (s *employeeService) GetEmployeeWorkdays(employeeID uint) ([]types.EmployeeWorkdayDTO, error) {
	workdays, err := s.repo.GetEmployeeWorkdaysByEmployeeID(employeeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employee workdays", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.EmployeeWorkdayDTO, len(workdays))
	for i, workday := range workdays {
		dtos[i] = *types.MapToEmployeeWorkdayDTO(&workday)
	}

	return dtos, nil
}

func (s *employeeService) GetAllWarehouseEmployees(warehouseID uint) ([]types.EmployeeAccountDTO, error) {
	if warehouseID == 0 {
		return nil, utils.WrapError("invalid warehouse ID", types.ErrValidation)
	}

	warehouseEmployees, err := s.repo.GetAllWarehouseEmployees(warehouseID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all warehouse employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(warehouseEmployees))
	for i, employee := range warehouseEmployees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) GetAllStoreEmployees(storeID uint) ([]types.EmployeeAccountDTO, error) {
	if storeID == 0 {
		return nil, utils.WrapError("invalid store ID", types.ErrValidation)
	}

	storeEmployees, err := s.repo.GetAllStoreEmployees(storeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(storeEmployees))
	for i, storeEmployee := range storeEmployees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&storeEmployee)
	}

	return dtos, nil
}

func (s *employeeService) GetAllRegionEmployees(regionID uint) ([]types.EmployeeAccountDTO, error) {
	if regionID == 0 {
		return nil, utils.WrapError("invalid region ID", types.ErrValidation)
	}

	storeEmployees, err := s.repo.GetAllRegionEmployees(regionID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all region employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(storeEmployees))
	for i, storeEmployee := range storeEmployees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&storeEmployee)
	}

	return dtos, nil
}

func (s *employeeService) GetAllFranchiseeEmployees(franchiseeID uint) ([]types.EmployeeAccountDTO, error) {
	if franchiseeID == 0 {
		return nil, utils.WrapError("invalid franchisee ID", types.ErrValidation)
	}

	storeEmployees, err := s.repo.GetAllFranchiseeEmployees(franchiseeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all franchisee employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(storeEmployees))
	for i, storeEmployee := range storeEmployees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&storeEmployee)
	}

	return dtos, nil
}

func (s *employeeService) GetAllAdminEmployees() ([]types.EmployeeAccountDTO, error) {
	adminEmployees, err := s.repo.GetAllAdminEmployees()
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all admin employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(adminEmployees))
	for i, employee := range adminEmployees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&employee)
	}

	return dtos, nil
}
