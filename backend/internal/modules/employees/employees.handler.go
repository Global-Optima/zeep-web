package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
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
}

func NewEmployeeHandler(service EmployeeService, auditService audit.AuditService, franchiseeService franchisees.FranchiseeService) *EmployeeHandler {
	return &EmployeeHandler{
		service:           service,
		auditService:      auditService,
		franchiseeService: franchiseeService,
	}
}

func (h *EmployeeHandler) DeleteStoreEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete warehouse employee: employee not found")
		return
	}

	err = h.service.DeleteTypedEmployee(uint(employeeID), storeID, data.StoreEmployeeType)
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete warehouse employee")
		return
	}

	action := types.DeleteStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, storeID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, "warehouse employee deleted successfully")
}

func (h *EmployeeHandler) DeleteWarehouseEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete warehouse employee: employee not found")
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

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, "warehouse employee deleted successfully")
}

func (h *EmployeeHandler) DeleteFranchiseeEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete franchisee employee: employee not found")
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

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, "franchisee employee deleted successfully")
}

func (h *EmployeeHandler) DeleteRegionManager(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	regionID, errH := contexts.GetRegionId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete region manager: employee not found")
		return
	}

	err = h.service.DeleteTypedEmployee(uint(employeeID), regionID, data.WarehouseRegionManagerEmployeeType)
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete region manager")
		return
	}

	action := types.DeleteRegionManagerEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, regionID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, "region manager deleted successfully")
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
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, "failed to fetch store id from context"+errH.Error(), errH.Status())
		return
	}

	ok, errH := h.FranchiseeStoreCheck(c, storeID)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}
	if !ok {
		utils.SendErrorWithStatus(c, "access denied", http.StatusForbidden)
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
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

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessCreatedResponse(c, "store employee created successfully")
}

func (h *EmployeeHandler) CreateWarehouseEmployee(c *gin.Context) {
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, "failed to fetch warehouse id from context: "+errH.Error(), errH.Status())
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
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

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessCreatedResponse(c, "warehouse employee created successfully")
}

func (h *EmployeeHandler) CreateFranchiseeEmployee(c *gin.Context) {
	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
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

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessCreatedResponse(c, "franchisee employee created successfully")
}

func (h *EmployeeHandler) CreateRegionManager(c *gin.Context) {
	regionID, errH := contexts.GetRegionId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	id, err := h.service.CreateRegionManager(regionID, &input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create region employee")
		return
	}

	action := types.CreateRegionManagerAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, regionID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessCreatedResponse(c, "region employee created successfully")
}

func (h *EmployeeHandler) GetStoreEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	storeID, errH := contexts.GetStoreId(c)
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

	warehouseID, errH := contexts.GetWarehouseId(c)
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

func (h *EmployeeHandler) GetRegionManagerByID(c *gin.Context) {
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

	employee, err := h.service.GetRegionManagerByID(uint(id), warehouseID)
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
	storeID, errH := contexts.GetStoreId(c)
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
	warehouseID, errH := contexts.GetWarehouseId(c)
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

func (h *EmployeeHandler) GetRegionManagers(c *gin.Context) {
	var filter types.EmployeesFilter
	regionID, errH := contexts.GetRegionId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.RegionManager{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetRegionManagers(regionID, &filter)
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

	storeID, errH := contexts.GetStoreId(c)
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

	err = h.service.UpdateStoreEmployee(uint(id), storeID, &input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update store employee")
		return
	}

	action := types.UpdateStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: storeEmployee.FirstName + " " + storeEmployee.LastName,
	}, &input, storeID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "employee updated successfully"})
}

func (h *EmployeeHandler) UpdateWarehouseEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	warehouseID, errH := contexts.GetWarehouseId(c)
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

	err = h.service.UpdateWarehouseEmployee(uint(id), warehouseID, &input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update warehouse employee")
		return
	}

	action := types.UpdateWarehouseEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: warehouseEmployee.FirstName + " " + warehouseEmployee.LastName,
	}, &input, warehouseID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "employee updated successfully"})
}

func (h *EmployeeHandler) UpdateFranchiseeEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	franchiseeID, errH := contexts.GetFranchiseeId(c)
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

	if err := h.service.UpdateFranchiseeEmployee(uint(id), franchiseeID, &input); err != nil {
		utils.SendInternalServerError(c, "failed to update franchisee employee")
		return
	}

	action := types.UpdateFranchiseeEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: franchiseeEmployee.FirstName + " " + franchiseeEmployee.LastName,
	}, &input, franchiseeID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

	utils.SendSuccessResponse(c, gin.H{"message": "franchisee employee updated successfully"})
}

func (h *EmployeeHandler) UpdateRegionManager(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	regionID, errH := contexts.GetRegionId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateRegionManagerEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, "invalid input")
		return
	}

	regionManagerEmployee, err := h.service.GetRegionManagerByID(uint(id), regionID)
	if err != nil {
		utils.SendBadRequestError(c, "failed to update region manager: employee not found")
		return
	}

	if err := h.service.UpdateRegionManager(uint(id), regionID, &input); err != nil {
		utils.SendInternalServerError(c, "failed to update region manager")
		return
	}

	action := types.UpdateRegionManagerEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: regionManagerEmployee.FirstName + " " + regionManagerEmployee.LastName,
	}, &input, regionID)

	_ = h.auditService.RecordEmployeeAction(c, &action)

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

func (h *EmployeeHandler) FranchiseeStoreCheck(c *gin.Context, storeID uint) (bool, *handlerErrors.HandlerError) {
	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		return false, errH
	}

	ok, err := h.franchiseeService.IsFranchiseeStore(franchiseeID, storeID)
	if err != nil {
		return false, errH
	}
	return ok, nil

}
