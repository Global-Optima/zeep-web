package employees

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
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
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	var input types.UpdatePasswordDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdatePassword(uint(id), &input)
	if err != nil {
		if err.Error() == "incorrect old password" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, "passwords mismatch")
			return
		}
		utils.SendInternalServerError(c, "failed to update password")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "password updated successfully"})
}

func (h *EmployeeHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve roles")
		return
	}

	utils.SendSuccessResponse(c, roles)
}

func (h *EmployeeHandler) GetCurrentEmployee(c *gin.Context) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		utils.SendErrorWithStatus(c, "failed to retrieve claims", http.StatusUnauthorized)
		return
	}

	employee, err := h.service.GetEmployeeByID(claims.EmployeeClaimsData.ID)
	if err != nil {
		utils.SendErrorWithStatus(c, "failed to retrieve employee", http.StatusUnauthorized)
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) ReassignEmployeeType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	var dto types.ReassignEmployeeTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		utils.SendBadRequestError(c, "failed to reassign employee type: employee not found")
		return
	}

	err = h.service.ReassignEmployeeType(uint(id), &dto)
	if err != nil {
		utils.SendInternalServerError(c, "failed to reassign employee type")
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

func (h *EmployeeHandler) CreateEmployeeWorkday(c *gin.Context) {
	var workday types.CreateEmployeeWorkdayDTO
	if err := c.ShouldBindJSON(&workday); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	_, err := h.service.CreateEmployeeWorkDay(&workday)
	if err != nil {
		utils.SendInternalServerError(c, "failed to create workday")
		return
	}

	utils.SendMessageWithStatus(c, "workday created successfully", http.StatusCreated)
}

func (h *EmployeeHandler) GetEmployeeWorkday(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	workday, err := h.service.GetEmployeeWorkday(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve workday")
		return
	}

	utils.SendSuccessResponse(c, workday)
}

func (h *EmployeeHandler) GetEmployeeWorkdays(c *gin.Context) {
	var employeeID uint

	id, err := strconv.ParseUint(c.Query("employeeId"), 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}
	employeeID = uint(id)

	workdays, err := h.service.GetEmployeeWorkdays(employeeID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve workdays")
		return
	}

	utils.SendSuccessResponse(c, workdays)
}

func (h *EmployeeHandler) GetMyWorkdays(c *gin.Context) {
	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		utils.SendErrorWithStatus(c, "failed to retrieve employee context", http.StatusUnauthorized)
		return
	}

	workdays, err := h.service.GetEmployeeWorkdays(claims.EmployeeClaimsData.ID)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve workdays")
		return
	}

	utils.SendSuccessResponse(c, workdays)
}

func (h *EmployeeHandler) UpdateEmployeeWorkday(c *gin.Context) {
	idParam := c.Param("id")
	workdayID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid workday ID")
		return
	}

	var updateWorkday types.UpdateEmployeeWorkdayDTO
	if err := c.ShouldBindJSON(&updateWorkday); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateEmployeeWorkday(uint(workdayID), &updateWorkday)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update workday")
		return
	}

	utils.SendMessageWithStatus(c, "workday updated successfully", http.StatusOK)
}

func (h *EmployeeHandler) DeleteEmployeeWorkday(c *gin.Context) {
	idParam := c.Param("id")
	workdayID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid workday ID")
		return
	}

	err = h.service.DeleteEmployeeWorkday(uint(workdayID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete workday")
		return
	}

	utils.SendMessageWithStatus(c, "workday deleted successfully", http.StatusOK)
}

func (h *EmployeeHandler) GetWarehouseAccounts(c *gin.Context) {
	warehouseIdStr := c.Param("id")
	warehouseID, err := strconv.ParseUint(warehouseIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid warehouse ID")
		return
	}

	warehouseEmployees, err := h.service.GetAllWarehouseEmployees(uint(warehouseID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employees")
		return
	}

	utils.SendSuccessResponse(c, warehouseEmployees)
}

func (h *EmployeeHandler) GetStoreAccounts(c *gin.Context) {
	storeIdStr := c.Param("id")
	storeID, err := strconv.ParseUint(storeIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	storeEmployees, err := h.service.GetAllStoreEmployees(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employees")
		return
	}

	utils.SendSuccessResponse(c, storeEmployees)
}

func (h *EmployeeHandler) GetRegionAccounts(c *gin.Context) {
	regionIdStr := c.Param("id")
	regionID, err := strconv.ParseUint(regionIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid region ID")
		return
	}

	regionEmployees, err := h.service.GetAllRegionEmployees(uint(regionID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve region employees")
		return
	}

	utils.SendSuccessResponse(c, regionEmployees)
}

func (h *EmployeeHandler) GetFranchiseeAccounts(c *gin.Context) {
	franchiseeIdStr := c.Param("id")
	franchiseeID, err := strconv.ParseUint(franchiseeIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid franchisee ID")
		return
	}

	franchiseeEmployees, err := h.service.GetAllFranchiseeEmployees(uint(franchiseeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve franchisee employees")
		return
	}

	utils.SendSuccessResponse(c, franchiseeEmployees)
}

func (h *EmployeeHandler) GetAdminAccounts(c *gin.Context) {
	adminEmployees, err := h.service.GetAllAdminEmployees()
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve admin employees")
		return
	}

	utils.SendSuccessResponse(c, adminEmployees)
}
