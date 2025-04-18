package middleware

import (
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"github.com/gin-gonic/gin"
)

var LocaleKey = "locale"

var AllowedLanguages = map[data.LanguageCode]struct{}{
	data.LanguageCodeEN: {},
	data.LanguageCodeRU: {},
	data.LanguageCodeKK: {},
}

func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var locale data.LanguageCode

		if q := c.Query("lang"); q != "" {
			q = normalizeLang(q)
			if isAllowed(data.LanguageCode(q)) {
				locale = data.LanguageCode(q)
			}
		}

		if locale == "" {
			header := c.GetHeader("Accept-Language")
			if header != "" {
				part := strings.SplitN(header, ",", 2)[0]
				part = strings.SplitN(part, ";", 2)[0]
				part = normalizeLang(part)
				sanitized := data.LanguageCode(part)
				if isAllowed(sanitized) {
					locale = sanitized
				}
			}
		}

		contexts.SetLocaleCtx(c, locale)
		c.Header("Content-Language", locale.String())
		c.Next()
	}
}

// normalizeLang trims whitespace and lowercases the code.
// You could expand this to strip region subtags, e.g. "en-US" → "en".
func normalizeLang(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func isAllowed(lang data.LanguageCode) bool {
	_, ok := AllowedLanguages[lang]
	return ok
}
