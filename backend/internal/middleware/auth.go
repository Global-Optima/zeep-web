package middleware

import (
	"net/http"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/employeeToken"
	authTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

func EmployeeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, err := authTypes.ExtractEmployeeAccessTokenAndValidate(c)
		if err != nil {
			zapLogger.Warn("missing or invalid token")
			utils.SendErrorWithStatus(c, "missing or invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		contexts.SetEmployeeCtx(c, claims)
		c.Next()
	}
}

func EmployeeSessionComparison(employeeTokenRepo employeeToken.EmployeeTokenRepository, employeesRepo employees.EmployeeRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		// Retrieve employee claims from context (set by EmployeeAuth middleware).
		claims, err := contexts.GetEmployeeClaimsFromCtx(c)
		if err != nil {
			zapLogger.Warn("no employee context found")
			utils.SendErrorWithStatus(c, "session invalid", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Extract the raw token from request (header or cookie).
		tokenStr, err := authTypes.ExtractToken(c, authTypes.ACCESS_TOKEN_HEADER, authTypes.EMPLOYEE_ACCESS_TOKEN_COOKIE_KEY)
		if err != nil {
			zapLogger.Warn("access token not found in request")
			utils.SendErrorWithStatus(c, "missing access token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Retrieve the stored session token (hashed) from the database.
		currentToken, err := employeeTokenRepo.GetTokenByEmployeeID(claims.EmployeeClaimsData.ID)
		if err != nil {
			zapLogger.Warnf("failed to retrieve session token from DB: %v", err)
			utils.SendErrorWithStatus(c, "session invalid", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Check if the stored token has expired.
		if currentToken.ExpiresAt.Before(time.Now()) {
			zapLogger.Warn("session token expired")
			utils.SendErrorWithStatus(c, "session expired", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Retrieve the current employee record from the database.
		currentEmployee, err := employeesRepo.GetEmployeeByID(claims.EmployeeClaimsData.ID)
		if err != nil {
			zapLogger.Warnf("failed to retrieve current employee record: %v", err)
			utils.SendErrorWithStatus(c, "session invalid", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Compare critical fields:
		// - Role: if the employee’s role has changed, then the token is outdated.
		// - WorkplaceID and EmployeeType: if the employee’s assignment has changed, then the token is outdated.
		currentClaims, err := authTypes.MapEmployeeToClaimsData(currentEmployee)
		if err != nil {
			zapLogger.Warnf("failed to map current employee claims: %v", err)
			utils.SendErrorWithStatus(c, "session invalid", http.StatusUnauthorized)
			c.Abort()
			return
		}
		if currentClaims.Role != claims.Role ||
			currentClaims.WorkplaceID != claims.WorkplaceID ||
			currentClaims.EmployeeType != claims.EmployeeType {
			zapLogger.Warn("employee privileges or assignment have changed; token outdated")
			utils.SendErrorWithStatus(c, "session outdated, please re-login", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Finally, compare the provided token with the stored hashed token using SHA-256.
		if err := utils.CompareTokenSHA256(currentToken.HashedToken, tokenStr); err != nil {
			zapLogger.Warn("token mismatch: provided token does not match current session token")
			utils.SendErrorWithStatus(c, "session token mismatch", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

func CustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		zapLogger := logger.GetZapSugaredLogger()

		claims, err := authTypes.ExtractCustomerAccessTokenAndValidate(c)
		if err != nil {
			zapLogger.Warn("missing or invalid token")
			utils.SendErrorWithStatus(c, "missing or invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		contexts.SetCustomerCtx(c, claims)
		c.Next()
	}
}
