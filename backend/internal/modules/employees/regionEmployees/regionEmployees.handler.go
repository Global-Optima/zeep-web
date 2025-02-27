package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees/types"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RegionEmployeeHandler struct {
	service         RegionEmployeeService
	employeeService employees.EmployeeService
	auditService    audit.AuditService
	regionService   regions.RegionService
}

func NewRegionEmployeeHandler(service RegionEmployeeService, employeeService employees.EmployeeService, regionService regions.RegionService, auditService audit.AuditService) *RegionEmployeeHandler {
	return &RegionEmployeeHandler{
		service:         service,
		employeeService: employeeService,
		regionService:   regionService,
		auditService:    auditService,
	}
}

func (h *RegionEmployeeHandler) DeleteRegionEmployee(c *gin.Context) {
	regionEmployeeIDParam := c.Param("id")
	regionEmployeeID, err := strconv.ParseUint(regionEmployeeIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400RegionEmployee)
		return
	}

	regionID, role, errH := contexts.GetRegionIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	regionEmployee, err := h.service.GetRegionEmployeeByID(uint(regionEmployeeID), regionID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionEmployeeGet)
		return
	}

	if !data.CanManageRole(role, regionEmployee.Role) {
		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(regionEmployeeID), regionEmployee.Region.ID, data.RegionEmployeeType)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionEmployeeDelete)
		return
	}

	action := types.DeleteRegionEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(regionEmployeeID),
		Name: regionEmployee.FirstName + " " + regionEmployee.LastName,
	}, struct{}{}, regionEmployee.Region.ID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200RegionEmployeeDelete)
}

func (h *RegionEmployeeHandler) CreateRegionEmployee(c *gin.Context) {
	regionID, role, errH := contexts.GetRegionIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	if regionID == nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400RegionEmployee)
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

	id, err := h.service.CreateRegionEmployee(*regionID, &input)
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrAlreadyExists):
			localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
			return
		case errors.Is(err, moduleErrors.ErrValidation):
			localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
			return
		default:
			localization.SendLocalizedResponseWithKey(c, types.Response500RegionEmployeeCreate)
			return
		}
	}

	action := types.CreateRegionEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, *regionID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201RegionEmployee)
}

func (h *RegionEmployeeHandler) GetRegionEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400RegionEmployee)
		return
	}

	warehouseID, errH := contexts.GetRegionId(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetRegionEmployeeByID(uint(id), warehouseID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionEmployeeGet)
		return
	}

	if employee == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
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

	if regionID == nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400RegionEmployee)
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.RegionEmployee{})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400RegionEmployee)
		return
	}

	regionEmployees, err := h.service.GetRegionEmployees(*regionID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionEmployeeGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, regionEmployees, filter.Pagination)
}

func (h *RegionEmployeeHandler) UpdateRegionEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400RegionEmployee)
		return
	}

	regionID, role, errH := contexts.GetRegionIdWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateRegionEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	regionEmployee, err := h.service.GetRegionEmployeeByID(uint(id), regionID)
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	if !data.CanManageRole(role, regionEmployee.Role) {
		localization.SendLocalizedResponseWithStatus(c, http.StatusForbidden)
		return
	}

	if err := h.service.UpdateRegionEmployee(uint(id), regionID, &input, role); err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionEmployeeUpdate)
		return
	}

	action := types.UpdateRegionEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: regionEmployee.FirstName + " " + regionEmployee.LastName,
	}, &input, regionEmployee.Region.ID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200RegionEmployeeUpdate)
}

func (h *RegionEmployeeHandler) GetRegionAccounts(c *gin.Context) {
	regionIdStr := c.Param("id")
	regionID, err := strconv.ParseUint(regionIdStr, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400RegionEmployee)
		return
	}

	regionEmployees, err := h.service.GetAllRegionEmployees(uint(regionID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500RegionEmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, regionEmployees)
}
