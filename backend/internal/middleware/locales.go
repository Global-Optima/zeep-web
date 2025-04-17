package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var AllowedLanguages = map[string]struct{}{
	"en": {},
	"ru": {},
	"kk": {},
}

const LocaleKey = "locale"

func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var locale string

		if q := c.Query("lang"); q != "" {
			q = normalizeLang(q)
			if isAllowed(q) {
				locale = q
			}
		}

		if locale == "" {
			header := c.GetHeader("Accept-Language")
			if header != "" {
				part := strings.SplitN(header, ",", 2)[0]
				part = strings.SplitN(part, ";", 2)[0]
				part = normalizeLang(part)
				if isAllowed(part) {
					locale = part
				}
			}
		}

		if locale == "" {
			locale = "ru"
		}

		c.Set(LocaleKey, locale)
		c.Header("Content-Language", locale)
		c.Next()
	}
}

// normalizeLang trims whitespace and lowercases the code.
// You could expand this to strip region subtags, e.g. "en-US" → "en".
func normalizeLang(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func isAllowed(lang string) bool {
	_, ok := AllowedLanguages[lang]
	return ok
}
