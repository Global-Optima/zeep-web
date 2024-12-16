package employees

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	EMPLOYEE_TOKEN_COOKIE_KEY = "EMPLOYEE_TOKEN"
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
		utils.SendBadRequestError(c, err.Error())
		return
	}

	employee, err := h.service.CreateEmployee(input)
	if err != nil {
		if err.Error() == "invalid email format" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, err.Error())
			return
		}
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	if employee == nil {
		utils.SendErrorWithStatus(c, "employee not found", http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	queryParams, err := types.ParseEmployeeQueryParams(c.Request.URL.Query())
	if err != nil {
		utils.SendBadRequestError(c, err.Error())
		utils.SendBadRequestError(c, err.Error())
		return
	}

	employees, err := h.service.GetEmployees(*queryParams)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, employees)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.SendBadRequestError(c, "invalid employee ID")
		return
	}

	var input types.UpdateEmployeeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	err = h.service.UpdateEmployee(uint(id), input)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "employee updated successfully"})
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
		utils.SendInternalServerError(c, err.Error())
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
		utils.SendBadRequestError(c, err.Error())
		return
	}

	err = h.service.UpdatePassword(uint(id), input)
	if err != nil {
		if err.Error() == "incorrect old password" || err.Error() == "password validation failed" {
			utils.SendBadRequestError(c, err.Error())
			return
		}
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "password updated successfully"})
}

func (h *EmployeeHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, roles)
}

func (h *EmployeeHandler) EmployeeLogin(c *gin.Context) {
	var input types.LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	token, err := h.service.EmployeeLogin(input.Email, input.Password)
	if err != nil {
		utils.SendErrorWithStatus(c, "invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := &utils.EmployeeClaims{}
	if err := utils.ValidateJWT(token, claims); err != nil {
		utils.SendInternalServerError(c, "failed to validate token")
		return
	}

	utils.SetCookie(c, EMPLOYEE_TOKEN_COOKIE_KEY, token, utils.CookieExpiration)

	utils.SendSuccessResponse(c, gin.H{"message": "login successful", "token": token})
}

func (h *EmployeeHandler) GetCurrentEmployee(c *gin.Context) {
	token, err := c.Cookie(EMPLOYEE_TOKEN_COOKIE_KEY)

	if err != nil {
		utils.SendErrorWithStatus(c, "authentication token missing", http.StatusUnauthorized)
		return
	}

	claims := &utils.EmployeeClaims{}
	if err := utils.ValidateJWT(token, claims); err != nil {
		utils.SendErrorWithStatus(c, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	employee, err := h.service.GetEmployeeByID(claims.ID)
	if err != nil {
		print(err)
		utils.SendInternalServerError(c, "failed to fetch employee details")
		return
	}

	print(employee)
	utils.SendSuccessResponse(c, employee)
}

func (h *EmployeeHandler) EmployeeLogout(c *gin.Context) {
	token, err := utils.GetCookie(c, EMPLOYEE_TOKEN_COOKIE_KEY)
	if err != nil {
		// Token not found in cookie
		utils.SendErrorWithStatus(c, "no token found", http.StatusUnauthorized)
		return
	}

	claims := &utils.EmployeeClaims{}
	if err := utils.ValidateJWT(token, claims); err != nil {
		utils.SendErrorWithStatus(c, "invalid token", http.StatusUnauthorized)
		return
	}

	utils.ClearCookie(c, EMPLOYEE_TOKEN_COOKIE_KEY)

	utils.SendSuccessResponse(c, gin.H{"message": "logout successful"})
}
