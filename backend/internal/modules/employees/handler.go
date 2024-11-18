package employees

import (
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

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var input types.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.service.CreateEmployee(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) GetEmployeesByStore(c *gin.Context) {
	storeIDParam := c.Query("storeId")
	roleIDParam := c.Query("roleId")
	limit, offset := utils.ParsePaginationParams(c)

	storeID, err := strconv.ParseUint(storeIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store ID"})
		return
	}

	var roleID *uint
	if roleIDParam != "" {
		parsedRoleID, err := strconv.ParseUint(roleIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
			return
		}
		roleIDUint := uint(parsedRoleID)
		roleID = &roleIDUint
	}

	employees, err := h.service.GetEmployeesByStore(uint(storeID), roleID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if employee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	var input types.UpdateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateEmployee(uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee updated successfully"})
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	err = h.service.DeleteEmployee(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee deleted successfully"})
}

func (h *EmployeeHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}
