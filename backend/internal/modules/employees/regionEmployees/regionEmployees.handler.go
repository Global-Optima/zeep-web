package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RegionEmployeeHandler struct {
	service           RegionEmployeeService
	employeeService   employees.EmployeeService
	auditService      audit.AuditService
	franchiseeService franchisees.FranchiseeService
	regionService     regions.RegionService
}

func NewRegionEmployeeHandler(service RegionEmployeeService, employeeService employees.EmployeeService, auditService audit.AuditService, franchiseeService franchisees.FranchiseeService, regionService regions.RegionService) *RegionEmployeeHandler {
	return &RegionEmployeeHandler{
		service:           service,
		employeeService:   employeeService,
		auditService:      auditService,
		franchiseeService: franchiseeService,
		regionService:     regionService,
	}
}

func (h *RegionEmployeeHandler) DeleteRegionEmployee(c *gin.Context) {
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

	employee, err := h.employeeService.GetEmployeeByID(uint(employeeID))
	if err != nil {
		utils.SendBadRequestError(c, "failed to delete region employee: employee not found")
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(employeeID), regionID, data.RegionEmployeeType)
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

func (h *RegionEmployeeHandler) CreateRegionEmployee(c *gin.Context) {
	regionID, role, errH := contexts.GetRegionIdWithRole(c)
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

func (h *RegionEmployeeHandler) GetRegionEmployeeByID(c *gin.Context) {
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

func (h *RegionEmployeeHandler) GetRegionEmployees(c *gin.Context) {
	var filter employeesTypes.EmployeesFilter
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

	regionEmployees, err := h.service.GetRegionEmployees(regionID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve region managers")
		return
	}

	utils.SendSuccessResponseWithPagination(c, regionEmployees, filter.Pagination)
}

func (h *RegionEmployeeHandler) UpdateRegionEmployee(c *gin.Context) {
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
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
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
