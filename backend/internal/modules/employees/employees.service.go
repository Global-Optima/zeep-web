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
	CreateStoreEmployee(storeID uint, input *types.CreateEmployeeDTO) (uint, error)
	CreateWarehouseEmployee(warehouseID uint, input *types.CreateEmployeeDTO) (uint, error)
	CreateRegionManager(regionID uint, input *types.CreateEmployeeDTO) (uint, error)
	CreateFranchiseeEmployee(franchiseeID uint, input *types.CreateEmployeeDTO) (uint, error)

	GetStoreEmployees(storeID uint, filter *types.EmployeesFilter) ([]types.StoreEmployeeDTO, error)
	GetWarehouseEmployees(warehouseID uint, filter *types.EmployeesFilter) ([]types.WarehouseEmployeeDTO, error)
	GetFranchiseeEmployees(franchiseeID uint, filter *types.EmployeesFilter) ([]types.FranchiseeEmployeeDTO, error)
	GetRegionManagers(regionID uint, filter *types.EmployeesFilter) ([]types.RegionManagerDTO, error)
	GetAdminEmployees(filter *types.EmployeesFilter) ([]types.AdminEmployeeDTO, error)

	GetEmployeeByID(id uint) (*types.EmployeeDTO, error)
	GetStoreEmployeeByID(id, storeID uint) (*types.StoreEmployeeDTO, error)
	GetWarehouseEmployeeByID(id, warehouseID uint) (*types.WarehouseEmployeeDTO, error)
	GetFranchiseeEmployeeByID(id, franchiseeID uint) (*types.FranchiseeEmployeeDTO, error)
	GetRegionManagerByID(id, regionID uint) (*types.RegionManagerDTO, error)
	GetAdminByID(id uint) (*types.AdminEmployeeDTO, error)

	UpdateFranchiseeEmployee(id, franchiseeID uint, input *types.UpdateFranchiseeEmployeeDTO) error
	UpdateStoreEmployee(id, storeID uint, input *types.UpdateStoreEmployeeDTO) error
	UpdateWarehouseEmployee(id, warehouseID uint, input *types.UpdateWarehouseEmployeeDTO) error
	UpdateRegionManager(id, regionID uint, input *types.UpdateRegionManagerEmployeeDTO) error

	DeleteEmployee(employeeID uint) error
	UpdatePassword(employeeID uint, input *types.UpdatePasswordDTO) error
	GetAllRoles() ([]types.EmployeeTypeRoles, error)

	CreateEmployeeWorkDay(dto *types.CreateEmployeeWorkdayDTO) (uint, error)
	GetEmployeeWorkday(workdayID uint) (*types.EmployeeWorkdayDTO, error)
	GetEmployeeWorkdays(employeeID uint) ([]types.EmployeeWorkdayDTO, error)
	UpdateEmployeeWorkday(workdayID uint, dto *types.UpdateEmployeeWorkdayDTO) error
	DeleteEmployeeWorkday(workdayID uint) error

	GetAllStoreEmployees(storeID uint) ([]types.EmployeeAccountDTO, error)
	GetAllWarehouseEmployees(warehouseID uint) ([]types.EmployeeAccountDTO, error)
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

func (s *employeeService) CreateStoreEmployee(storeID uint, input *types.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToStoreEmployee(storeID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create store employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

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

func (s *employeeService) CreateWarehouseEmployee(warehouseID uint, input *types.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToWarehouseEmployee(warehouseID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create warehouse employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

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
		wrappedErr := utils.WrapError("failed to create warehouse employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *employeeService) CreateFranchiseeEmployee(franchiseeID uint, input *types.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToFranchiseeEmployee(franchiseeID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create franchisee employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

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
		wrappedErr := utils.WrapError("failed to create franchisee employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *employeeService) CreateRegionManager(regionID uint, input *types.CreateEmployeeDTO) (uint, error) {
	employee, err := types.CreateToRegionManager(regionID, input)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create region manager employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

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
		wrappedErr := utils.WrapError("failed to create region manager employee", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	return id, nil
}

func (s *employeeService) GetStoreEmployees(storeID uint, filter *types.EmployeesFilter) ([]types.StoreEmployeeDTO, error) {
	employees, err := s.repo.GetStoreEmployees(storeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.StoreEmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToStoreEmployeeDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) GetWarehouseEmployees(warehouseID uint, filter *types.EmployeesFilter) ([]types.WarehouseEmployeeDTO, error) {
	employees, err := s.repo.GetWarehouseEmployees(warehouseID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.WarehouseEmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToWarehouseEmployeeDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) GetFranchiseeEmployees(franchiseeID uint, filter *types.EmployeesFilter) ([]types.FranchiseeEmployeeDTO, error) {
	employees, err := s.repo.GetFranchiseeEmployees(franchiseeID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve franchisee employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.FranchiseeEmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToFranchiseeEmployeeDTO(&employee)
	}
	return dtos, nil
}

func (s *employeeService) GetRegionManagers(regionID uint, filter *types.EmployeesFilter) ([]types.RegionManagerDTO, error) {
	employees, err := s.repo.GetRegionManagers(regionID, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve region managers", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.RegionManagerDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToRegionManagerDTO(&employee)
	}
	return dtos, nil
}

func (s *employeeService) GetAdminEmployees(filter *types.EmployeesFilter) ([]types.AdminEmployeeDTO, error) {
	employees, err := s.repo.GetAdminEmployees(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve admin employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	dtos := make([]types.AdminEmployeeDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToAdminEmployeeDTO(&employee)
	}
	return dtos, nil
}

func (s *employeeService) GetEmployeeByID(id uint) (*types.EmployeeDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid store employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(id)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve employee for store employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToEmployeeDTO(employee), nil
}

func (s *employeeService) GetStoreEmployeeByID(id, storeID uint) (*types.StoreEmployeeDTO, error) {
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

func (s *employeeService) GetWarehouseEmployeeByID(id, warehouseID uint) (*types.WarehouseEmployeeDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid warehouse employee ID")
	}

	employee, err := s.repo.GetWarehouseEmployeeByID(id, warehouseID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve warehouse employee with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToWarehouseEmployeeDTO(employee), nil
}

func (s *employeeService) GetFranchiseeEmployeeByID(id, franchiseeID uint) (*types.FranchiseeEmployeeDTO, error) {
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

	return types.MapToFranchiseeEmployeeDTO(employee), nil
}

func (s *employeeService) GetRegionManagerByID(id, regionID uint) (*types.RegionManagerDTO, error) {
	if id == 0 {
		return nil, errors.New("invalid region manager ID")
	}

	employee, err := s.repo.GetRegionManagerByID(id, regionID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to retrieve region manager with ID = %d: %w", id, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if employee == nil {
		return nil, errors.New("employee not found")
	}

	return types.MapToRegionManagerDTO(employee), nil
}

func (s *employeeService) GetAdminByID(id uint) (*types.AdminEmployeeDTO, error) {
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

func (s *employeeService) UpdateFranchiseeEmployee(id, franchiseeID uint, input *types.UpdateFranchiseeEmployeeDTO) error {
	updateFields, err := types.FranchiseeEmployeeUpdateFields(input)
	if err != nil {
		return err
	}
	return s.repo.UpdateFranchiseeEmployee(id, franchiseeID, updateFields)
}

func (s *employeeService) UpdateRegionManager(id, regionID uint, input *types.UpdateRegionManagerEmployeeDTO) error {
	updateFields, err := types.RegionManagerEmployeeUpdateFields(input)
	if err != nil {
		return err
	}
	return s.repo.UpdateRegionManager(id, regionID, updateFields)
}

func (s *employeeService) UpdateStoreEmployee(id, storeID uint, input *types.UpdateStoreEmployeeDTO) error {
	updateFields, err := types.StoreEmployeeUpdateFields(input)
	if err != nil {
		return err
	}
	return s.repo.UpdateStoreEmployee(id, storeID, updateFields)
}

func (s *employeeService) UpdateWarehouseEmployee(id, warehouseID uint, input *types.UpdateWarehouseEmployeeDTO) error {
	updateFields, err := types.WarehouseEmployeeUpdateFields(input)
	if err != nil {
		return err
	}
	return s.repo.UpdateWarehouseEmployee(id, warehouseID, updateFields)
}

func (s *employeeService) DeleteEmployee(employeeID uint) error {
	if employeeID == 0 {
		return errors.New("invalid employee ID")
	}

	employee, err := s.repo.GetEmployeeByID(employeeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve employee", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	if err := s.repo.DeleteEmployeeById(employeeID, employee.Type); err != nil {
		wrappedErr := fmt.Errorf("failed to delete employee with ID = %d: %w", employeeID, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

// TODO move to Recovery module or leave for ADMIN only
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

	employee.HashedPassword = hashedPassword
	if err := s.repo.UpdateEmployee(employeeID, employee); err != nil {
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

func (s *employeeService) CreateEmployeeWorkDay(dto *types.CreateEmployeeWorkdayDTO) (uint, error) {
	errMsg := "failed to create employee work day"

	workday, err := types.ValidateEmployeeWorkday(dto)
	if err != nil {
		wrappedErr := utils.WrapError(errMsg, err)
		s.logger.Error(wrappedErr)
		return 0, err
	}

	existingWorkday, err := s.repo.GetEmployeeWorkdayByEmployeeAndDay(workday.EmployeeID, workday.Day)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr := utils.WrapError(errMsg, err)
		s.logger.Error(wrappedErr)
		return 0, err
	}
	if existingWorkday != nil {
		err = fmt.Errorf("%w: not unique workday for employeeID %d in %v ", types.ErrWorkdayAlreadyExists, dto.EmployeeID, dto.Day)
		wrappedErr := utils.WrapError(errMsg, err)
		s.logger.Error(wrappedErr)
		return 0, err
	}

	id, err := s.repo.CreateEmployeeWorkday(workday)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create employee workday", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	return id, nil
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

func (s *employeeService) UpdateEmployeeWorkday(workdayID uint, dto *types.UpdateEmployeeWorkdayDTO) error {
	workday, err := types.WorkdaysUpdateFields(dto)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update employee workday", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	err = s.repo.UpdateEmployeeWorkdayById(workdayID, workday)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update employee workday", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *employeeService) DeleteEmployeeWorkday(workdayID uint) error {
	err := s.repo.DeleteEmployeeWorkday(workdayID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to delete employee workday", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *employeeService) GetAllStoreEmployees(storeID uint) ([]types.EmployeeAccountDTO, error) {
	if storeID == 0 {
		return nil, utils.WrapError("invalid store ID", types.ErrValidation)
	}

	employees, err := s.repo.GetAllStoreEmployees(storeID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all store employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&employee)
	}

	return dtos, nil
}
func (s *employeeService) GetAllWarehouseEmployees(warehouseID uint) ([]types.EmployeeAccountDTO, error) {
	if warehouseID == 0 {
		return nil, utils.WrapError("invalid warehouse ID", types.ErrValidation)
	}

	employees, err := s.repo.GetAllWarehouseEmployees(warehouseID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all warehouse employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&employee)
	}

	return dtos, nil
}

func (s *employeeService) GetAllAdminEmployees() ([]types.EmployeeAccountDTO, error) {
	employees, err := s.repo.GetAllAdminEmployees()
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve all admin employees", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAccountDTO, len(employees))
	for i, employee := range employees {
		dtos[i] = *types.MapToEmployeeAccountDTO(&employee)
	}

	return dtos, nil
}
