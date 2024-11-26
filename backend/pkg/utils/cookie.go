package utils

import (
	"net/http"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	CookieExpiration = 8 * time.Hour
	CookiePath       = "/"
)

func SetCookie(c *gin.Context, name, value string, expiration time.Duration) {
	cfg := config.GetConfig()

	secure := true
	domain := cfg.Server.ClientURL
	sameSite := http.SameSiteNoneMode

	if cfg.IsDevelopment {
		secure = false
		sameSite = http.SameSiteLaxMode
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     CookiePath,
		Domain:   domain,
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
	domain := cfg.Server.ClientURL
	sameSite := http.SameSiteNoneMode

	if cfg.IsDevelopment {
		secure = false
		sameSite = http.SameSiteLaxMode
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     CookiePath,
		Domain:   domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}

	http.SetCookie(c.Writer, cookie)
}
