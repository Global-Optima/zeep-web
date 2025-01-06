package employees

import (
	"fmt"
	"net/http"
	"strconv"


	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service EmployeeService
}

func NewEmployeeHandler(service EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	err = h.service.DeleteEmployee(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete employee")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "employee deleted successfully"})
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

	var storeEmployee *types.StoreEmployeeDTO
	var warehouseEmployee *types.WarehouseEmployeeDTO

	switch claims.EmployeeType {
	case data.StoreEmployeeType:
		storeEmployee, err = h.service.GetStoreEmployeeByID(claims.EmployeeClaimsData.ID)
		if err != nil {
			print(err)
			utils.SendInternalServerError(c, "failed to fetch employee details")
			return
		}
		utils.SendSuccessResponse(c, storeEmployee)
		return
	case data.WarehouseEmployeeType:
		warehouseEmployee, err = h.service.GetWarehouseEmployeeByID(claims.EmployeeClaimsData.ID)
		if err != nil {
			print(err)
			utils.SendInternalServerError(c, "failed to fetch employee details")
			return
		}
		utils.SendSuccessResponse(c, warehouseEmployee)
		return
	}

	utils.SendBadRequestError(c, fmt.Sprintf("invalid employee type: %v", claims.EmployeeType))
}

func (h *EmployeeHandler) CreateStoreEmployee(c *gin.Context) {
	var input types.CreateStoreEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	employee, err := h.service.CreateStoreEmployee(&input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create store employee")
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) CreateWarehouseEmployee(c *gin.Context) {
	var input types.CreateWarehouseEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	employee, err := h.service.CreateWarehouseEmployee(&input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
			return
		}
		utils.SendInternalServerError(c, "failed to create warehouse employee")
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetStoreEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	employee, err := h.service.GetStoreEmployeeByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employee details")
		return
	}

	if employee == nil {
		utils.SendErrorWithStatus(c, "employee not found", http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) GetWarehouseEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	employee, err := h.service.GetWarehouseEmployeeByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employee details")
		return
	}

	if employee == nil {
		utils.SendErrorWithStatus(c, "employee not found", http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) GetStoreEmployees(c *gin.Context) {
	var filter types.GetStoreEmployeesFilter
	storeID, errH := contexts.GetStoreId(c)
	if errH != nil && !errors.Is(errH, contexts.ErrEmptyStoreID) {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Employee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetStoreEmployees(storeID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, employees, filter.Pagination)
}

func (h *EmployeeHandler) GetWarehouseEmployees(c *gin.Context) {
	var filter types.GetWarehouseEmployeesFilter
	warehouseID, errH := contexts.GetWarehouseId(c)
	if errH != nil && !errors.Is(errH, contexts.ErrEmptyWarehouseID) {
		utils.SendErrorWithStatus(c, errH.Error(), errH.Status())
		return
	}

	err := utils.ParseQueryWithBaseFilter(c, &filter, &data.Employee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetWarehouseEmployees(warehouseID, &filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employees")
		return
	}

	utils.SendSuccessResponseWithPagination(c, employees, filter.Pagination)
}

func (h *EmployeeHandler) UpdateStoreEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	var input types.UpdateStoreEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateStoreEmployee(uint(id), &input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update store employee")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "employee updated successfully"})
}

func (h *EmployeeHandler) UpdateWarehouseEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	var input types.UpdateWarehouseEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.UpdateWarehouseEmployee(uint(id), &input)
	if err != nil {
		utils.SendInternalServerError(c, "failed to update warehouse employee")
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "employee updated successfully"})
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

	claims, err := contexts.GetEmployeeClaimsFromCtx(c)
	if err != nil {
		utils.SendErrorWithStatus(c, "failed to retrieve employee context", http.StatusUnauthorized)
		return
	}

	if claims.Role != data.RoleAdmin && claims.Role != data.RoleDirector {
		employeeID = claims.EmployeeClaimsData.ID
	} else {
		id, err := strconv.ParseUint(c.Query("employeeId"), 10, 64)
		if err != nil {
			utils.SendBadRequestError(c, "invalid employee ID")
			return
		}
		employeeID = uint(id)
	}

	workdays, err := h.service.GetEmployeeWorkdays(employeeID)
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

	var updateWorkday types.UpdateEmployeeWorkdayDTO
	if err := c.ShouldBindJSON(&updateWorkday); err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_JSON)
		return
	}

	err = h.service.DeleteEmployeeWorkday(uint(workdayID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to delete workday")
		return
	}

	utils.SendMessageWithStatus(c, "workday deleted successfully", http.StatusOK)
}

func (h *EmployeeHandler) GetStoreAccounts(c *gin.Context) {
	storeIdStr := c.Param("id")
	storeID, err := strconv.ParseUint(storeIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid store ID")
		return
	}

	employees, err := h.service.GetAllStoreEmployees(uint(storeID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
}

func (h *EmployeeHandler) GetWarehouseAccounts(c *gin.Context) {
	warehouseIdStr := c.Param("id")
	warehouseID, err := strconv.ParseUint(warehouseIdStr, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid warehouse ID")
		return
	}

	employees, err := h.service.GetAllWarehouseEmployees(uint(warehouseID))
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
}

func (h *EmployeeHandler) GetAdminAccounts(c *gin.Context) {
	employees, err := h.service.GetAllAdminEmployees()
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve admin employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
}
