package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	CookieExpiration = 8 * time.Hour
	CookiePath       = "/"
)

// Detect if Client URL is using HTTP (no HTTPS)
func isInsecureMode(clientURL string) bool {
	return strings.HasPrefix(clientURL, "http://") // Check if URL starts with HTTP
}

func SetCookie(c *gin.Context, name, value string, expiration time.Duration) {
	cfg := config.GetConfig()

	secure := true
	sameSite := http.SameSiteNoneMode

	// If running in development mode OR using an IP without HTTPS, make cookies less strict
	if cfg.IsDevelopment || isInsecureMode(cfg.Server.ClientURL) {
		secure = false
		sameSite = http.SameSiteLaxMode
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     CookiePath,
		Expires:  time.Now().Add(expiration),
		MaxAge:   int(expiration.Seconds()),
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}

	http.SetCookie(c.Writer, cookie)
}

func GetCookie(c *gin.Context, name string) (string, error) {
	return c.Cookie(name)
}

func ClearCookie(c *gin.Context, name string) {
	cfg := config.GetConfig()

	secure := true
	sameSite := http.SameSiteNoneMode

	if cfg.IsDevelopment || isInsecureMode(cfg.Server.ClientURL) {
		secure = false
		sameSite = http.SameSiteLaxMode
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     CookiePath,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}

	http.SetCookie(c.Writer, cookie)
}
