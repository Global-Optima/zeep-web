package employees

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/adminEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AdminEmployeeHandler struct {
	service         AdminEmployeeService
	employeeService employees.EmployeeService
	auditService    audit.AuditService
}

func NewAdminEmployeeHandler(service AdminEmployeeService, employeeService employees.EmployeeService, auditService audit.AuditService) *AdminEmployeeHandler {
	return &AdminEmployeeHandler{
		service:         service,
		employeeService: employeeService,
		auditService:    auditService,
	}
}

func (h *AdminEmployeeHandler) CreateAdminEmployee(c *gin.Context) {
	var input employeesTypes.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	id, err := h.service.CreateAdminEmployee(&input)
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrAlreadyExists):
			localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
			return
		case errors.Is(err, moduleErrors.ErrValidation):
			localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500AdminEmployeeCreate)
			return
		}
	}

	action := types.CreateAdminEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201AdminEmployee)
}

func (h *AdminEmployeeHandler) GetAdminEmployees(c *gin.Context) {
	var filter employeesTypes.EmployeesFilter

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.AdminEmployee{})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	adminEmployees, err := h.service.GetAdminEmployees(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdminEmployeeGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, adminEmployees, filter.Pagination)
}

func (h *AdminEmployeeHandler) GetAdminEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400AdminEmployee)
		return
	}

	employee, err := h.service.GetAdminEmployeeByID(uint(id))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdminEmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *AdminEmployeeHandler) GetAdminAccounts(c *gin.Context) {
	adminEmployees, err := h.service.GetAllAdminEmployees()
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500AdminEmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, adminEmployees)
}
