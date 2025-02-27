package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type FranchiseeEmployeeHandler struct {
	service           FranchiseeEmployeeService
	employeeService   employees.EmployeeService
	auditService      audit.AuditService
	franchiseeService franchisees.FranchiseeService
}

func NewFranchiseeEmployeeHandler(service FranchiseeEmployeeService, employeeService employees.EmployeeService, franchiseeService franchisees.FranchiseeService, auditService audit.AuditService) *FranchiseeEmployeeHandler {
	return &FranchiseeEmployeeHandler{
		service:           service,
		employeeService:   employeeService,
		franchiseeService: franchiseeService,
		auditService:      auditService,
	}
}

func (h *FranchiseeEmployeeHandler) DeleteFranchiseeEmployee(c *gin.Context) {
	franchiseeEmployeeIDParam := c.Param("id")
	franchiseeEmployeeID, err := strconv.ParseUint(franchiseeEmployeeIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400FranchiseeEmployee)
		return
	}

	franchiseeID, role, errH := contexts.GetFranchiseeIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	franchiseeEmployee, err := h.service.GetFranchiseeEmployeeByID(uint(franchiseeEmployeeID), franchiseeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeDelete)
		return
	}

	if !data.CanManageRole(role, franchiseeEmployee.Role) {
		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(franchiseeEmployeeID), franchiseeEmployee.Franchisee.ID, data.FranchiseeEmployeeType)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeDelete)
		return
	}

	action := types.DeleteFranchiseeEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(franchiseeEmployeeID),
		Name: franchiseeEmployee.FirstName + " " + franchiseeEmployee.LastName,
	}, struct{}{}, franchiseeEmployee.Franchisee.ID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200FranchiseeEmployeeDelete)
}

func (h *FranchiseeEmployeeHandler) CreateFranchiseeEmployee(c *gin.Context) {
	franchiseeID, role, errH := contexts.GetFranchiseeIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input employeesTypes.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	if !data.CanManageRole(role, input.Role) {
		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		return
	}

	if franchiseeID == nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400FranchiseeEmployee)
		return
	}

	id, err := h.service.CreateFranchiseeEmployee(*franchiseeID, &input)
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrAlreadyExists):
			localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
			return
		case errors.Is(err, moduleErrors.ErrValidation):
			localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeCreate)
			return
		}
	}

	action := types.CreateFranchiseeEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, *franchiseeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201FranchiseeEmployee)
}

func (h *FranchiseeEmployeeHandler) GetFranchiseeEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400FranchiseeEmployee)
		return
	}

	franchiseeID, errH := contexts.GetFranchiseeId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetFranchiseeEmployeeByID(uint(id), franchiseeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeGet)
		return
	}

	if employee == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
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

	if franchiseeID == nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400FranchiseeEmployee)
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.FranchiseeEmployee{})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	franchiseeEmployees, err := h.service.GetFranchiseeEmployees(*franchiseeID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, franchiseeEmployees, filter.Pagination)
}

func (h *FranchiseeEmployeeHandler) UpdateFranchiseeEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400FranchiseeEmployee)
		return
	}

	franchiseeID, role, errH := contexts.GetFranchiseeIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateFranchiseeEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	franchiseeEmployee, err := h.service.GetFranchiseeEmployeeByID(uint(id), franchiseeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeGet)
		return
	}
	if !data.CanManageRole(role, franchiseeEmployee.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	if err := h.service.UpdateFranchiseeEmployee(uint(id), franchiseeID, &input, role); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeUpdate)
		return
	}

	action := types.UpdateFranchiseeEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: franchiseeEmployee.FirstName + " " + franchiseeEmployee.LastName,
	}, &input, franchiseeEmployee.Franchisee.ID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200FranchiseeEmployeeUpdate)
}

func (h *FranchiseeEmployeeHandler) GetFranchiseeAccounts(c *gin.Context) {
	franchiseeIdStr := c.Param("id")
	franchiseeID, err := strconv.ParseUint(franchiseeIdStr, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400FranchiseeEmployee)
		return
	}

	franchiseeEmployees, err := h.service.GetAllFranchiseeEmployees(uint(franchiseeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500FranchiseeEmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, franchiseeEmployees)
}
