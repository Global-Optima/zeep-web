package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	employeesTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/warehouseEmployees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type WarehouseEmployeeHandler struct {
	service         WarehouseEmployeeService
	employeeService employees.EmployeeService
	regionService   regions.RegionService
	auditService    audit.AuditService
}

func NewWarehouseEmployeeHandler(service WarehouseEmployeeService, employeeService employees.EmployeeService, regionService regions.RegionService, auditService audit.AuditService) *WarehouseEmployeeHandler {
	return &WarehouseEmployeeHandler{
		service:         service,
		employeeService: employeeService,
		regionService:   regionService,
		auditService:    auditService,
	}
}

func (h *WarehouseEmployeeHandler) DeleteWarehouseEmployee(c *gin.Context) {
	employeeIDParam := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400WarehouseEmployee)
		return
	}

	warehouseID, role, errH := h.regionService.CheckRegionWarehouseWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.employeeService.GetEmployeeByID(uint(employeeID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeDelete)
		return
	}

	if !data.CanManageRole(role, employee.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.employeeService.DeleteTypedEmployee(uint(employeeID), warehouseID, data.WarehouseEmployeeType)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeDelete)
		return
	}

	action := types.DeleteWarehouseEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(employeeID),
		Name: employee.FirstName + " " + employee.LastName,
	}, struct{}{}, warehouseID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseEmployeeDelete)
}

func (h *WarehouseEmployeeHandler) CreateWarehouseEmployee(c *gin.Context) {
	warehouseID, role, errH := h.regionService.CheckRegionWarehouseWithRole(c)
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

	id, err := h.service.CreateWarehouseEmployee(warehouseID, &input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			return
		}
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeCreate)
		return
	}

	action := types.CreateWarehouseEmployeeAuditFactory(&data.BaseDetails{
		ID:   id,
		Name: input.FirstName + " " + input.LastName,
	}, &input, warehouseID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response201WarehouseEmployee)
}

func (h *WarehouseEmployeeHandler) GetWarehouseEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400WarehouseEmployee)
		return
	}

	warehouseID, errH := h.regionService.CheckRegionWarehouse(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	employee, err := h.service.GetWarehouseEmployeeByID(uint(id), warehouseID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeGet)
		return
	}

	if employee == nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *WarehouseEmployeeHandler) GetWarehouseEmployees(c *gin.Context) {
	var filter employeesTypes.EmployeesFilter
	warehouseID, errH := h.regionService.CheckRegionWarehouse(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.WarehouseEmployee{})
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingQuery)
		return
	}

	warehouseEmployees, err := h.service.GetWarehouseEmployees(warehouseID, &filter)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeGet)
		return
	}

	utils.SendSuccessResponseWithPagination(c, warehouseEmployees, filter.Pagination)
}

func (h *WarehouseEmployeeHandler) UpdateWarehouseEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400WarehouseEmployee)
		return
	}

	warehouseID, role, errH := h.regionService.CheckRegionWarehouseWithRole(c)
	if errH != nil {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	var input types.UpdateWarehouseEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		return
	}

	warehouseEmployee, err := h.service.GetWarehouseEmployeeByID(uint(id), warehouseID)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeGet)
		return
	}
	if !data.CanManageRole(role, warehouseEmployee.Role) {
		utils.SendErrorWithStatus(c, employeesTypes.ErrNotAllowedToManageTheRole.Error(), http.StatusForbidden)
		return
	}

	err = h.service.UpdateWarehouseEmployee(uint(id), warehouseID, &input, role)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeUpdate)
		return
	}

	action := types.UpdateWarehouseEmployeeAuditFactory(&data.BaseDetails{
		ID:   uint(id),
		Name: warehouseEmployee.FirstName + " " + warehouseEmployee.LastName,
	}, &input, warehouseID)

	go func() {
		_ = h.auditService.RecordEmployeeAction(c, &action)
	}()

	localization.SendLocalizedResponseWithKey(c, types.Response200WarehouseEmployeeUpdate)
}

func (h *WarehouseEmployeeHandler) GetWarehouseAccounts(c *gin.Context) {
	warehouseIdStr := c.Param("id")
	warehouseID, err := strconv.ParseUint(warehouseIdStr, 10, 64)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response400WarehouseEmployee)
		return
	}

	warehouseEmployees, err := h.service.GetAllWarehouseEmployees(uint(warehouseID))
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, types.Response500WarehouseEmployeeGet)
		return
	}

	utils.SendSuccessResponse(c, warehouseEmployees)
}
