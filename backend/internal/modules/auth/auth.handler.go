package auth

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
	"net/http"

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

	tokenPair, err := h.service.CustomerLogin(input.Phone, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, moduleErrors.ErrValidation):
			utils.SendErrorWithStatus(c, "invalid credentials", http.StatusBadRequest)
			return
		default:
			utils.SendErrorWithStatus(c, "unexpected error", http.StatusInternalServerError)
			return
		}
	}

	claims := &types.CustomerClaims{}
	if err := types.ValidateCustomerJWT(tokenPair.AccessToken, claims, types.TokenAccess); err != nil {
		utils.SendInternalServerError(c, "failed to validate token")
		return
	}

	cfg := config.GetConfig()

	utils.SetCookie(c, types.CUSTOMER_ACCESS_TOKEN_COOKIE_KEY, tokenPair.AccessToken, cfg.JWT.CustomerAccessTokenTTL)
	utils.SetCookie(c, types.CUSTOMER_REFRESH_TOKEN_COOKIE_KEY, tokenPair.RefreshToken, cfg.JWT.CustomerRefreshTokenTTL)

	utils.SendSuccessResponse(c, gin.H{
		"message": "login successful",
		"data": gin.H{
			"accessToken":  tokenPair.AccessToken,
			"refreshToken": tokenPair.RefreshToken,
		},
	})
}

func (h *AuthenticationHandler) CustomerRefresh(c *gin.Context) {
	token, err := types.ExtractToken(c, types.REFRESH_TOKEN_HEADER, types.CUSTOMER_REFRESH_TOKEN_COOKIE_KEY)
	if err != nil {
		// Token not found in cookie
		utils.SendErrorWithStatus(c, "no token found", http.StatusUnauthorized)
		return
	}

	tokenPair, err := h.service.CustomerRefreshTokens(token)
	if err != nil {
		utils.SendErrorWithStatus(c, "invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	cfg := config.GetConfig()

	utils.SetCookie(c, types.CUSTOMER_ACCESS_TOKEN_COOKIE_KEY, tokenPair.AccessToken, cfg.JWT.CustomerAccessTokenTTL)
	utils.SetCookie(c, types.CUSTOMER_REFRESH_TOKEN_COOKIE_KEY, tokenPair.RefreshToken, cfg.JWT.CustomerRefreshTokenTTL)

	utils.SendSuccessResponse(c, gin.H{
		"message": "refresh successful",
		"data": gin.H{
			"accessToken":  tokenPair.AccessToken,
			"refreshToken": tokenPair.RefreshToken,
		},
	})
}

func (h *AuthenticationHandler) CustomerLogout(c *gin.Context) {
	utils.ClearCookie(c, types.CUSTOMER_ACCESS_TOKEN_COOKIE_KEY)
	utils.ClearCookie(c, types.CUSTOMER_REFRESH_TOKEN_COOKIE_KEY)

	utils.SendSuccessResponse(c, gin.H{"message": "logout successful"})
}

func (h *AuthenticationHandler) EmployeeLogin(c *gin.Context) {
	var input types.EmployeeLoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendBadRequestError(c, err.Error())
		return
	}

	tokenPair, err := h.service.EmployeeLogin(input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrInvalidCredentials):
			utils.SendErrorWithStatus(c, "invalid credentials", http.StatusBadRequest)
			return
		default:
			utils.SendErrorWithStatus(c, "unexpected error", http.StatusInternalServerError)
			return
		}
	}

	claims := &types.EmployeeClaims{}
	if err := types.ValidateEmployeeJWT(tokenPair.AccessToken, claims, types.TokenAccess); err != nil {
		utils.SendInternalServerError(c, "failed to validate token ")
		return
	}

	cfg := config.GetConfig()

	utils.SetCookie(c, types.EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY, tokenPair.AccessToken, cfg.JWT.EmployeeAccessTokenTTL)
	utils.SetCookie(c, types.EMPLOYEE_REFRESH_TOKEN_COOKIE_KEY, tokenPair.RefreshToken, cfg.JWT.EmployeeRefreshTokenTTL)

	utils.SendSuccessResponse(c, gin.H{
		"message": "login successful",
		"data": gin.H{
			"accessToken":  tokenPair.AccessToken,
			"refreshToken": tokenPair.RefreshToken,
		},
	})
}

func (h *AuthenticationHandler) EmployeeRefresh(c *gin.Context) {
	refreshToken, err := types.ExtractToken(c, types.REFRESH_TOKEN_HEADER, types.EMPLOYEE_REFRESH_TOKEN_COOKIE_KEY)
	if err != nil {
		// Token not found in cookie
		utils.SendErrorWithStatus(c, "no token found", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := h.service.EmployeeRefreshAccessToken(refreshToken)
	if err != nil {
		utils.SendErrorWithStatus(c, "invalid token", http.StatusUnauthorized)
		return
	}

	cfg := config.GetConfig()

	utils.SetCookie(c, types.EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY, newAccessToken, cfg.JWT.CustomerAccessTokenTTL)

	utils.SendSuccessResponse(c, gin.H{
		"message": "refresh successful",
		"data": gin.H{
			"accessToken": newAccessToken,
		},
	})
}

func (h *AuthenticationHandler) EmployeeLogout(c *gin.Context) {
	utils.ClearCookie(c, types.EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY)
	utils.ClearCookie(c, types.EMPLOYEE_REFRESH_TOKEN_COOKIE_KEY)

	utils.SendSuccessResponse(c, gin.H{"message": "logout successful"})
}
