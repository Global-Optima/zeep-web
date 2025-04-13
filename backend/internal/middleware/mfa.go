package middleware

import (
	authTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RequireMFA() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _, err := authTypes.ExtractEmployeeSessionTokenAndValidate(c)
		if err != nil {
			logrus.Infof("invalid token: %v", err)
			utils.SendErrorWithStatus(c, "missing or invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if !claims.MFA {
			// Если не прошёл MFA — перенаправляем на страницу авторизации 2FA
			c.Redirect(http.StatusFound, "/employee/mfa")
			c.Abort()
			return
		}
		c.Next()
	}
}
