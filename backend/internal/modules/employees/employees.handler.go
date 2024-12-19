package employees

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"net/http"
	"strconv"

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
	case data.WarehouseEmployeeType:
		warehouseEmployee, err = h.service.GetWarehouseEmployeeByID(claims.EmployeeClaimsData.ID)
		if err != nil {
			print(err)
			utils.SendInternalServerError(c, "failed to fetch employee details")
			return
		}
		utils.SendSuccessResponse(c, warehouseEmployee)
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
	var filter *types.GetStoreEmployeesFilter

	err := utils.ParseQueryWithBaseFilter(c, filter, &data.Employee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetStoreEmployees(filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve store employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
}

func (h *EmployeeHandler) GetWarehouseEmployees(c *gin.Context) {
	var filter *types.GetWarehouseEmployeesFilter

	err := utils.ParseQueryWithBaseFilter(c, filter, &data.Employee{})
	if err != nil {
		utils.SendBadRequestError(c, utils.ERROR_MESSAGE_BINDING_QUERY)
		return
	}

	employees, err := h.service.GetWarehouseEmployees(filter)
	if err != nil {
		utils.SendInternalServerError(c, "failed to retrieve warehouse employees")
		return
	}

	utils.SendSuccessResponse(c, employees)
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
