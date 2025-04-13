package auth

import (
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	service AuthenticationService
}

func NewAuthenticationHandler(service AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{service: service}
}

func (h *AuthenticationHandler) CustomerRegister(c *gin.Context) {
	var input types.CustomerRegisterDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	_, err := h.service.CustomerRegister(&input)
	if err != nil {
		utils.SendInternalServerError(c, err.Error())
		return
	}

	utils.SendSuccessResponse(c, gin.H{"message": "new user registered"})
}

func (h *AuthenticationHandler) CustomerLogin(c *gin.Context) {
	var input types.CustomerLoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	token, err := h.service.CustomerLogin(input.Phone, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrValidation):
			localization.SendLocalizedResponseWithKey(c, types.Response400IncorrectCredentials)
			return
		default:
			utils.SendErrorWithStatus(c, "unexpected error", http.StatusInternalServerError)
			return
		}
	}

	claims := &types.CustomerClaims{}
	if err := types.ValidateCustomerToken(token.SessionToken, claims); err != nil {
		utils.SendInternalServerError(c, "failed to validate token")
		return
	}

	cfg := config.GetConfig()

	utils.SetCookie(c, types.CUSTOMER_SESSION_COOKIE_KEY, token.SessionToken, cfg.JWT.CustomerTokenTTL)

	utils.SendSuccessResponse(c, gin.H{
		"message": "login successful",
		"data": gin.H{
			"sessionToken": token.SessionToken,
		},
	})
}

func (h *AuthenticationHandler) CustomerLogout(c *gin.Context) {
	utils.ClearCookie(c, types.CUSTOMER_SESSION_COOKIE_KEY)

	utils.SendSuccessResponse(c, gin.H{"message": "logout successful"})
}

// EmployeeLogin godoc
// @Summary Employee login
// @Description Authenticates an employee and returns a session token
// @Tags employee-auth
// @Accept json
// @Produce json
// @Param input body types.EmployeeLoginDTO true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful, session token returned"
// @Failure 400 {object} map[string]interface{} "Invalid credentials or bad request"
// @Failure 403 {object} map[string]interface{} "Employee account inactive"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/employees/employee/login [post]
func (h *AuthenticationHandler) EmployeeLogin(c *gin.Context) {
	var input types.EmployeeLoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	token, err := h.service.EmployeeLogin(input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrInvalidCredentials):
			localization.SendLocalizedResponseWithKey(c, types.Response400IncorrectCredentials)
		case errors.Is(err, types.ErrInactiveEmployee):
			utils.SendErrorWithStatus(c, "inactive employee", http.StatusForbidden)
		default:
			utils.SendErrorWithStatus(c, "unexpected error", http.StatusInternalServerError)
		}
		return
	}

	var claims types.EmployeeClaims
	if err := types.ValidateEmployeeToken(token.SessionToken, &claims); err != nil {
		utils.SendInternalServerError(c, "failed to validate token")
		return
	}

	cfg := config.GetConfig()

	utils.SetCookie(c, types.EMPLOYEE_SESSION_COOKIE_KEY, token.SessionToken, cfg.JWT.EmployeeTokenTTL)

	utils.SendSuccessResponse(c, gin.H{
		"message": "login successful",
		"data": gin.H{
			"sessionToken": token.SessionToken,
		},
	})
}

// EmployeeLogout godoc
// @Summary Employee logout
// @Description Logs out the currently authenticated employee
// @Tags employee-auth
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Router /api/v1/auth/employees/logout [post]
func (h *AuthenticationHandler) EmployeeLogout(c *gin.Context) {
	c.Set(contexts.EMPLOYEE_CONTEXT, nil)

	utils.ClearCookie(c, types.EMPLOYEE_SESSION_COOKIE_KEY)
	utils.SendSuccessResponse(c, gin.H{"message": "logout successful"})
}

func (h *AuthenticationHandler) EmployeeMFAPass(c *gin.Context) {

	_, _, err := types.ExtractEmployeeSessionTokenAndValidate(c)
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, 403)
		return
	}

	img, err := media.GetImageWithFormFile(c)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageGettingImage)
		return
	}

	token, err := h.service.EmployeeFaceRecognitionPass(c, img)
	if err != nil {
		logrus.Info(err)
		localization.SendLocalizedResponseWithKey(c, types.Response400IncorrectCredentials)
		return
	}

	cfg := config.GetConfig()

	utils.SetCookie(c, types.EMPLOYEE_SESSION_COOKIE_KEY, token.SessionToken, cfg.JWT.EmployeeTokenTTL)

	utils.SendSuccessResponse(c, gin.H{
		"message": "login successful",
		"data": gin.H{
			"sessionToken": token.SessionToken,
		},
	})
}
