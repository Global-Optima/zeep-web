package contexts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/gin-gonic/gin"
)

var LocaleKey = "locale"

func GetLocaleFromCtx(c *gin.Context) data.LanguageCode {
	locale, ok := c.Get(LocaleKey)
	if !ok {
		return data.LanguageCodeRU
	}

	if lang, ok := locale.(data.LanguageCode); ok && lang.IsValid() {
		return lang
	}

	return data.LanguageCodeRU
}

func SetLocaleCtx(c *gin.Context, locale data.LanguageCode) {
	c.Set(LocaleKey, locale)
}
