package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type FranchiseeEmployeeHandler struct {
	service           FranchiseeEmployeeService
	employeeService   employees.EmployeeService
	auditService      audit.AuditService
	franchiseeService franchisees.FranchiseeService
	regionService     regions.RegionService
}

func NewFranchiseeEmployeeHandler(service FranchiseeEmployeeService, employeeService employees.EmployeeService, auditService audit.AuditService, franchiseeService franchisees.FranchiseeService, regionService regions.RegionService) *FranchiseeEmployeeHandler {
	return &FranchiseeEmployeeHandler{
		service:           service,
		employeeService:   employeeService,
		auditService:      auditService,
		franchiseeService: franchiseeService,
		regionService:     regionService,
	}
}

func (h *FranchiseeEmployeeHandler) DeleteFranchiseeEmployee(c *gin.Context) {
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

	employee, err := h.employeeService.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete franchisee employee: employee not found")
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(employeeID), franchiseeID, data.FranchiseeEmployeeType)
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

func (h *FranchiseeEmployeeHandler) CreateFranchiseeEmployee(c *gin.Context) {
	franchiseeID, role, errH := contexts.GetFranchiseeIdWithRole(c)
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

func (h *FranchiseeEmployeeHandler) GetFranchiseeEmployeeByID(c *gin.Context) {
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

func (h *FranchiseeEmployeeHandler) GetFranchiseeEmployees(c *gin.Context) {
	var filter employeesTypes.EmployeesFilter
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

	franchiseeEmployees, err := h.service.GetFranchiseeEmployees(franchiseeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisee employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, franchiseeEmployees, filter.Pagination)
}

func (h *FranchiseeEmployeeHandler) UpdateFranchiseeEmployee(c *gin.Context) {
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
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
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
