package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/adminEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AdminEmployeeHandler struct {
	service           AdminEmployeeService
	employeeService   employees.EmployeeService
	auditService      audit.AuditService
	franchiseeService franchisees.FranchiseeService
	regionService     regions.RegionService
}

func NewAdminEmployeeHandler(service AdminEmployeeService, employeeService employees.EmployeeService, auditService audit.AuditService, franchiseeService franchisees.FranchiseeService, regionService regions.RegionService) *AdminEmployeeHandler {
	return &AdminEmployeeHandler{
		service:           service,
		employeeService:   employeeService,
		auditService:      auditService,
		franchiseeService: franchiseeService,
		regionService:     regionService,
	}
}

func (h *AdminEmployeeHandler) CreateAdminEmployee(c *gin.Context) {
	var input employeesTypes.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	id, err := h.service.CreateAdminEmployee(&input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create franchisee employee")
		return
	}

	action := types.CreateAdminEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessCreatedResponse(c, "franchisee employee created successfully")
}

func (h *AdminEmployeeHandler) GetAdminEmployees(c *gin.Context) {
	var filter employeesTypes.EmployeesFilter

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.AdminEmployee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	adminEmployees, err := h.service.GetAdminEmployees(&filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve admin employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, adminEmployees, filter.Pagination)
}

func (h *AdminEmployeeHandler) GetAdminEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	employee, err := h.service.GetAdminEmployeeByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve admin employee")
		return
	}

	utils.SendSuccessResponse(c, employee)
}
