package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SanitizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.ContentType()

		switch contentType {
		case "application/json":
			processJSONRequest(c)
		case "multipart/form-data":
			processMultipartRequest(c)
		}

		c.Next()
	}
}

func processJSONRequest(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	for key, value := range requestData {
		if strVal, ok := value.(string); ok {
			sanitized, valid := utils.SanitizeString(strVal)
			if !valid {
				localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
				return
			}
			requestData[key] = sanitized
		}
	}

	sanitizedBody, _ := json.Marshal(requestData)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
}

func processMultipartRequest(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
		return
	}

	for key, values := range form.Value {
		for i, val := range values {
			sanitized, valid := utils.SanitizeString(val)
			if !valid {
				localization.SendLocalizedResponseWithStatus(c, http.StatusBadRequest)
				return
			}
			form.Value[key][i] = sanitized
		}
	}

	c.Request.MultipartForm = form
}
