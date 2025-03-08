package employees

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/pkg/errors"

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
	storeEmployeeIDParam := c.Param("id")
	storeEmployeeID, err := strconv.ParseUint(storeEmployeeIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400StoreEmployee)
		return
	}

	filter, errH := contexts.GetStoreContextFilter(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
		return
	}

	storeEmployee, err := h.service.GetStoreEmployeeByID(uint(storeEmployeeID), filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeDelete)
		return
	}

	if !data.CanManageRole(claims.Role, storeEmployee.Role) {
		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(storeEmployeeID), storeEmployee.Store.ID, data.StoreEmployeeType)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeDelete)
		return
	}

	action := types.DeleteStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(storeEmployeeID),
		Name: storeEmployee.FirstName + " " + storeEmployee.LastName,
	}, struct{}{}, storeEmployee.Store.ID)

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
		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		return
	}

	id, err := h.service.CreateStoreEmployee(storeID, &input)
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrAlreadyExists):
			localization.SendLocalizedResponseWithKey(c, types.Response409StoreEmployee)
			return
		case errors.Is(err, moduleErrors.ErrValidation):
			localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeCreate)
			return
		}
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

	filter, errH := contexts.GetStoreContextFilter(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetStoreEmployeeByID(uint(id), filter)
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

	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusUnauthorized)
		return
	}

	filter, errH := contexts.GetStoreContextFilter(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateStoreEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	storeEmployee, err := h.service.GetStoreEmployeeByID(uint(id), filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeGet)
		return
	}

	if !data.CanManageRole(claims.Role, storeEmployee.Role) {
		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		return
	}

	err = h.service.UpdateStoreEmployee(uint(id), filter, &input, claims.Role)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500StoreEmployeeUpdate)
		return
	}

	action := types.UpdateStoreEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: storeEmployee.FirstName + " " + storeEmployee.LastName,
	}, &input, storeEmployee.Store.ID)

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
