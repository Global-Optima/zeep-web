package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
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
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreEmployee)
		return
	}

	storeID, role, errH := h.franchiseeService.CheckFranchiseeStoreWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.employeeService.GetEmployeeByID(uint(employeeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeDelete)
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(employeeID), storeID, data.StoreEmployeeType)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeDelete)
		return
	}

	action := types.DeleteStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreEmployeeDelete)
}

func (h *StoreEmployeeHandler) CreateStoreEmployee(c *gin.Context) {
	storeID, role, errH := h.franchiseeService.CheckFranchiseeStoreWithRole(c)
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
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	id, err := h.service.CreateStoreEmployee(storeID, &input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeCreate)
		return
	}

	action := types.CreateStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201StoreEmployee)
}

func (h *StoreEmployeeHandler) GetStoreEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreEmployee)
		return
	}

	storeID, errH := h.franchiseeService.CheckFranchiseeStore(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetStoreEmployeeByID(uint(id), storeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeGet)
		return
	}

	if employee == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
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
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	storeEmployees, err := h.service.GetStoreEmployees(storeID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, storeEmployees, filter.Pagination)
}

func (h *StoreEmployeeHandler) UpdateStoreEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreEmployee)
		return
	}

	storeID, role, errH := h.franchiseeService.CheckFranchiseeStoreWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateStoreEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	storeEmployee, err := h.service.GetStoreEmployeeByID(uint(id), storeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeGet)
		return
	}

	err = h.service.UpdateStoreEmployee(uint(id), storeID, &input, role)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeUpdate)
		return
	}

	action := types.UpdateStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: storeEmployee.FirstName + " " + storeEmployee.LastName,
	}, &input, storeID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200StoreEmployeeUpdate)
}

func (h *StoreEmployeeHandler) GetStoreAccounts(c *gin.Context) {
	storeIdStr := c.Param("id")
	storeID, err := strconv.ParseUint(storeIdStr, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreEmployee)
		return
	}

	storeEmployees, err := h.service.GetAllStoreEmployees(uint(storeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, storeEmployees)
}
