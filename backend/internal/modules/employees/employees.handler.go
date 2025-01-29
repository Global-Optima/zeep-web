package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service           EmployeeService
	auditService      audit.AuditService
	franchiseeService franchisees.FranchiseeService
	regionService     regions.RegionService
}

func NewEmployeeHandler(service EmployeeService, auditService audit.AuditService, franchiseeService franchisees.FranchiseeService, regionService regions.RegionService) *EmployeeHandler {
	return &EmployeeHandler{
		service:           service,
		auditService:      auditService,
		franchiseeService: franchiseeService,
		regionService:     regionService,
	}
}

func (h *EmployeeHandler) DeleteStoreEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	storeID, role, errH := h.franchiseeService.CheckFranchiseeStoreWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete store employee: employee not found")
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.service.DeleteTypedEmployee(uint(employeeID), storeID, data.StoreEmployeeType)
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete store employee")
		return
	}

	action := types.DeleteStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, "store employee deleted successfully")
}

func (h *EmployeeHandler) DeleteWarehouseEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	warehouseID, role, errH := h.regionService.CheckRegionWarehouseWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete warehouse employee: employee not found")
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.service.DeleteTypedEmployee(uint(employeeID), warehouseID, data.WarehouseEmployeeType)
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete warehouse employee")
		return
	}

	action := types.DeleteWarehouseEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, warehouseID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, "warehouse employee deleted successfully")
}

func (h *EmployeeHandler) DeleteFranchiseeEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	franchiseeID, role, errH := contexts.GetFranchiseeIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete franchisee employee: employee not found")
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.service.DeleteTypedEmployee(uint(employeeID), franchiseeID, data.FranchiseeEmployeeType)
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete franchisee employee")
		return
	}

	action := types.DeleteFranchiseeEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, franchiseeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, "franchisee employee deleted successfully")
}

func (h *EmployeeHandler) DeleteRegionEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	regionID, role, errH := contexts.GetRegionIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete region employee: employee not found")
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.service.DeleteTypedEmployee(uint(employeeID), regionID, data.RegionEmployeeType)
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete region employee")
		return
	}

	action := types.DeleteRegionEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, regionID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, "region employee deleted successfully")
}

func (h *EmployeeHandler) UpdatePassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	var input types.UpdatePasswordDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdatePassword(uint(id), &input)
	if err != nil {
		if err.Error() == "incorrect old password" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, "passwords mismatch")
			return
		}
		utils.SendInternalServerError(c, "failed to update password")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "password updated successfully"})
}

func (h *EmployeeHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve roles")
		return
	}

	utils.SendSuccessResponse(c, roles)
}

func (h *EmployeeHandler) GetCurrentEmployee(c *gin.Context) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		utils.SendErrorWithStatus(c, "failed to retrieve claims", http.StatusUnauthorized)
		return
	}

	employee, err := h.service.GetEmployeeByID(claims.EmployeeClaimsData.ID)
	if err != nil {
		utils.SendErrorWithStatus(c, "failed to retrieve employee", http.StatusUnauthorized)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) CreateStoreEmployee(c *gin.Context) {
	storeID, role, errH := h.franchiseeService.CheckFranchiseeStoreWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if !data.CanManageRole(role, input.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	id, err := h.service.CreateStoreEmployee(storeID, &input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create store employee")
		return
	}

	action := types.CreateStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "store employee created successfully")
}

func (h *EmployeeHandler) CreateWarehouseEmployee(c *gin.Context) {
	warehouseID, role, errH := h.regionService.CheckRegionWarehouseWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	// Check if the role is manageable by the current user
	if !data.CanManageRole(role, input.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	id, err := h.service.CreateWarehouseEmployee(warehouseID, &input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create warehouse employee")
		return
	}

	action := types.CreateWarehouseEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, warehouseID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "warehouse employee created successfully")
}

func (h *EmployeeHandler) CreateFranchiseeEmployee(c *gin.Context) {
	franchiseeID, role, errH := contexts.GetFranchiseeIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if !data.CanManageRole(role, input.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	id, err := h.service.CreateFranchiseeEmployee(franchiseeID, &input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create franchisee employee")
		return
	}

	action := types.CreateFranchiseeEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, franchiseeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "franchisee employee created successfully")
}

func (h *EmployeeHandler) CreateRegionEmployee(c *gin.Context) {
	regionID, role, errH := contexts.GetRegionIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	// Check if the role is manageable by the current user
	if !data.CanManageRole(role, input.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	id, err := h.service.CreateRegionEmployee(regionID, &input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create region employee")
		return
	}

	action := types.CreateRegionEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, regionID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "region employee created successfully")
}

func (h *EmployeeHandler) GetStoreEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetStoreEmployeeByID(uint(id), storeID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employee details")
		return
	}

	if employee == nil {
		utils.SendErrorWithStatus(c, "employee not found", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetWarehouseEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	warehouseID, errH := h.regionService.CheckRegionWarehouse(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetWarehouseEmployeeByID(uint(id), warehouseID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employee details")
		return
	}

	if employee == nil {
		utils.SendErrorWithStatus(c, "employee not found", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetFranchiseeEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	storeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetFranchiseeEmployeeByID(uint(id), storeID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisee employee details")
		return
	}

	if employee == nil {
		utils.SendErrorWithStatus(c, "employee not found", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetRegionEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	warehouseID, errH := contexts.GetRegionId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetRegionEmployeeByID(uint(id), warehouseID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employee details")
		return
	}

	if employee == nil {
		utils.SendErrorWithStatus(c, "employee not found", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetStoreEmployees(c *gin.Context) {
	var filter types.EmployeesFilter
	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.StoreEmployee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetStoreEmployees(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, employees, filter.Pagination)
}

func (h *EmployeeHandler) GetWarehouseEmployees(c *gin.Context) {
	var filter types.EmployeesFilter
	warehouseID, errH := h.regionService.CheckRegionWarehouse(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.WarehouseEmployee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetWarehouseEmployees(warehouseID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, employees, filter.Pagination)
}

func (h *EmployeeHandler) GetFranchiseeEmployees(c *gin.Context) {
	var filter types.EmployeesFilter
	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.FranchiseeEmployee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetFranchiseeEmployees(franchiseeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisee employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, employees, filter.Pagination)
}

func (h *EmployeeHandler) GetRegionEmployees(c *gin.Context) {
	var filter types.EmployeesFilter
	regionID, errH := contexts.GetRegionId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.RegionEmployee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetRegionEmployees(regionID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve region managers")
		return
	}

	utils.SendSuccessResponseWithPagination(c, employees, filter.Pagination)
}

func (h *EmployeeHandler) GetAdminEmployees(c *gin.Context) {
	var filter types.EmployeesFilter

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.AdminEmployee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetAdminEmployees(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve admin employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, employees, filter.Pagination)
}

func (h *EmployeeHandler) UpdateStoreEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	storeID, role, errH := h.franchiseeService.CheckFranchiseeStoreWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateStoreEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	storeEmployee, err := h.service.GetStoreEmployeeByID(uint(id), storeID)
	if err != nil {
		utils.SendBadRequestError(c, "failed to update store employee: employee not found")
		return
	}
	if !data.CanManageRole(role, storeEmployee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.service.UpdateStoreEmployee(uint(id), storeID, &input, role)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update store employee")
		return
	}

	action := types.UpdateStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: storeEmployee.FirstName + " " + storeEmployee.LastName,
	}, &input, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "employee updated successfully"})
}

func (h *EmployeeHandler) UpdateWarehouseEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	warehouseID, role, errH := h.regionService.CheckRegionWarehouseWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateWarehouseEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	warehouseEmployee, err := h.service.GetWarehouseEmployeeByID(uint(id), warehouseID)
	if err != nil {
		utils.SendBadRequestError(c, "failed to update warehouse employee: employee not found")
		return
	}
	if !data.CanManageRole(role, warehouseEmployee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.service.UpdateWarehouseEmployee(uint(id), warehouseID, &input, role)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update warehouse employee")
		return
	}

	action := types.UpdateWarehouseEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: warehouseEmployee.FirstName + " " + warehouseEmployee.LastName,
	}, &input, warehouseID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "employee updated successfully"})
}

func (h *EmployeeHandler) UpdateFranchiseeEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	franchiseeID, role, errH := contexts.GetFranchiseeIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateFranchiseeEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}

	franchiseeEmployee, err := h.service.GetFranchiseeEmployeeByID(uint(id), franchiseeID)
	if err != nil {
		utils.SendBadRequestError(c, "failed to update franchisee employee: employee not found")
		return
	}
	if !data.CanManageRole(role, franchiseeEmployee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	if err := h.service.UpdateFranchiseeEmployee(uint(id), franchiseeID, &input, role); err != nil {
		utils.SendInternalServerError(c, "failed to update franchisee employee")
		return
	}

	action := types.UpdateFranchiseeEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: franchiseeEmployee.FirstName + " " + franchiseeEmployee.LastName,
	}, &input, franchiseeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "franchisee employee updated successfully"})
}

func (h *EmployeeHandler) UpdateRegionEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	regionID, role, errH := contexts.GetRegionIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateRegionEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}

	regionEmployee, err := h.service.GetRegionEmployeeByID(uint(id), regionID)
	if err != nil {
		utils.SendBadRequestError(c, "failed to update region manager: employee not found")
		return
	}
	if !data.CanManageRole(role, regionEmployee.Role) {
		utils.SendErrorWithStatus(c, types.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	if err := h.service.UpdateRegionEmployee(uint(id), regionID, &input, role); err != nil {
		utils.SendInternalServerError(c, "failed to update region manager")
		return
	}

	action := types.UpdateRegionEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: regionEmployee.FirstName + " " + regionEmployee.LastName,
	}, &input, regionID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "region manager updated successfully"})
}

func (h *EmployeeHandler) CreateEmployeeWorkday(c *gin.Context) {
	var workday types.CreateEmployeeWorkdayDTO
	if err := c.ShouldBindJSON(&workday); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateEmployeeWorkDay(&workday)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create workday")
		return
	}

	utils.SendMessageWithStatus(c, "workday created successfully", http.StatusCreated)
}

// TODO separate by types
func (h *EmployeeHandler) GetEmployeeWorkday(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	workday, err := h.service.GetEmployeeWorkday(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve workday")
		return
	}

	utils.SendSuccessResponse(c, workday)
}

func (h *EmployeeHandler) GetEmployeeWorkdays(c *gin.Context) {
	var employeeID uint

	id, err := strconv.ParseUint(c.Query("employeeId"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}
	employeeID = uint(id)

	workdays, err := h.service.GetEmployeeWorkdays(employeeID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve workdays")
		return
	}

	utils.SendSuccessResponse(c, workdays)
}

func (h *EmployeeHandler) GetMyWorkdays(c *gin.Context) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		utils.SendErrorWithStatus(c, "failed to retrieve employee context", http.StatusUnauthorized)
		return
	}

	workdays, err := h.service.GetEmployeeWorkdays(claims.EmployeeClaimsData.ID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve workdays")
		return
	}

	utils.SendSuccessResponse(c, workdays)
}

func (h *EmployeeHandler) UpdateEmployeeWorkday(c *gin.Context) {
	idParam := c.Param("id")
	workdayID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid workday ID")
		return
	}

	var updateWorkday types.UpdateEmployeeWorkdayDTO
	if err := c.ShouldBindJSON(&updateWorkday); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateEmployeeWorkday(uint(workdayID), &updateWorkday)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update workday")
		return
	}

	utils.SendMessageWithStatus(c, "workday updated successfully", http.StatusOK)
}

func (h *EmployeeHandler) DeleteEmployeeWorkday(c *gin.Context) {
	idParam := c.Param("id")
	workdayID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid workday ID")
		return
	}

	err = h.service.DeleteEmployeeWorkday(uint(workdayID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete workday")
		return
	}

	utils.SendMessageWithStatus(c, "workday deleted successfully", http.StatusOK)
}

func (h *EmployeeHandler) GetStoreAccounts(c *gin.Context) {
	storeIdStr := c.Param("id")
	storeID, err := strconv.ParseUint(storeIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	employees, err := h.service.GetAllStoreEmployees(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
}

func (h *EmployeeHandler) GetWarehouseAccounts(c *gin.Context) {
	warehouseIdStr := c.Param("id")
	warehouseID, err := strconv.ParseUint(warehouseIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid warehouse ID")
		return
	}

	employees, err := h.service.GetAllWarehouseEmployees(uint(warehouseID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
}

func (h *EmployeeHandler) GetAdminAccounts(c *gin.Context) {
	employees, err := h.service.GetAllAdminEmployees()
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve admin employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
}
