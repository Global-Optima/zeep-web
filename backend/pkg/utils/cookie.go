package utils

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	CookieExpiration = 8 * time.Hour // General cookie expiration time
	CookiePath       = "/"           // Cookie path
	CookieSecure     = true          // Secure cookies for HTTPS
	CookieHttpOnly   = true          // Restrict cookie access to HTTP(S)
	SameSiteStrict   = "Strict"      // SameSite policy for strict cross-site behavior
	SameSiteLax      = "Lax"         // Alternative SameSite policy if needed
)

func SetCookie(c *gin.Context, name, value string, expiration time.Duration) {
	cfg := config.GetConfig()

	c.SetCookie(
		name,
		value,
		int(expiration.Seconds()),
		CookiePath,
		cfg.Server.ClientURL,
		CookieSecure,
		CookieHttpOnly,
	)
}

func GetCookie(c *gin.Context, name string) (string, error) {
	return c.Cookie(name)
}

func ClearCookie(c *gin.Context, name string) {
	cfg := config.GetConfig()

	c.SetCookie(
		name,
		"",
		-1,
		CookiePath,
		cfg.Server.ClientURL,
		CookieSecure,
		CookieHttpOnly,
	)
}
