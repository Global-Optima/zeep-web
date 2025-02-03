package employees

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StoreEmployeeHandler struct {
	service           StoreEmployeeService
	employeeService   employees.EmployeeService
	auditService      audit.AuditService
	franchiseeService franchisees.FranchiseeService
}

func NewStoreEmployeeHandler(service StoreEmployeeService, employeeService employees.EmployeeService, franchiseeService franchisees.FranchiseeService, auditService audit.AuditService) *StoreEmployeeHandler {
	return &StoreEmployeeHandler{
		service:           service,
		employeeService:   employeeService,
		franchiseeService: franchiseeService,
		auditService:      auditService,
	}
}

func (h *StoreEmployeeHandler) DeleteStoreEmployee(c *gin.Context) {
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

	employee, err := h.employeeService.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete store employee: employee not found")
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(employeeID), storeID, data.StoreEmployeeType)
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

func (h *StoreEmployeeHandler) CreateStoreEmployee(c *gin.Context) {
	storeID, role, errH := h.franchiseeService.CheckFranchiseeStoreWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input employeesTypes.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	if !data.CanManageRole(role, input.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
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

func (h *StoreEmployeeHandler) GetStoreEmployeeByID(c *gin.Context) {
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

func (h *StoreEmployeeHandler) GetStoreEmployees(c *gin.Context) {
	var filter employeesTypes.EmployeesFilter
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

	storeEmployees, err := h.service.GetStoreEmployees(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, storeEmployees, filter.Pagination)
}

func (h *StoreEmployeeHandler) UpdateStoreEmployee(c *gin.Context) {
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
