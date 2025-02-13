package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"net/http"
	"strconv"

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

func (h *EmployeeHandler) UpdatePassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Employee)
		return
	}

	var input types.UpdatePasswordDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	err = h.service.UpdatePassword(uint(id), &input)
	if err != nil {
		if err.Error() == "incorrect old password" || err.Error() == "password validation failed" {
			localization.SendLocalizedResponseWithKey(c, types.Response400Employee)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500EmployeeUpdatePassword)
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "password updated successfully"})
}

func (h *EmployeeHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, roles)
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Employee)
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500EmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	var filter types.EmployeesFilter

	if err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Employee{}); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	employee, err := h.service.GetEmployees(&filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500EmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetCurrentEmployee(c *gin.Context) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response401Employee)
		return
	}

	employee, err := h.service.GetEmployeeByID(claims.EmployeeClaimsData.ID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response401Employee)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) ReassignEmployeeType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Employee)
		return
	}

	var dto types.ReassignEmployeeTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Employee)
		return
	}

	err = h.service.ReassignEmployeeType(uint(id), &dto)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500EmployeeReassignType)
		return
	}

	action := types.UpdateEmployeeAuditFactory(
		&data.BaseDetails{
			ID:   uint(id),
			Name: employee.FirstName + employee.LastName,
		},
		&dto)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	utils.SendSuccessResponse(c, gin.H{"message": "reassign employee type successfully"})
}

func (h *EmployeeHandler) GetEmployeeWorkday(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Employee)
		return
	}

	workday, err := h.service.GetEmployeeWorkday(uint(id))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500EmployeeGetWorkday)
		return
	}

	utils.SendSuccessResponse(c, workday)
}

func (h *EmployeeHandler) GetEmployeeWorkdays(c *gin.Context) {
	var employeeID uint

	id, err := strconv.ParseUint(c.Query("employeeId"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400Employee)
		return
	}
	employeeID = uint(id)

	workdays, err := h.service.GetEmployeeWorkdays(employeeID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500EmployeeGetWorkdays)
		return
	}

	utils.SendSuccessResponse(c, workdays)
}

func (h *EmployeeHandler) GetMyWorkdays(c *gin.Context) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response401Employee)
		return
	}

	workdays, err := h.service.GetEmployeeWorkdays(claims.EmployeeClaimsData.ID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500EmployeeGetWorkdays)
		return
	}

	utils.SendSuccessResponse(c, workdays)
}
