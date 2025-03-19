package middleware

import (
	"bytes"
	"encoding/json"
	"io"

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
	var requestData interface{}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		c.Abort()
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &requestData); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		c.Abort()
		return
	}

	requestData = sanitizeRecursive(requestData)

	sanitizedBody, _ := json.Marshal(requestData)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
}

func sanitizeRecursive(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeRecursive(value)
		}
		return v
	case []interface{}:
		for i, item := range v {
			v[i] = sanitizeRecursive(item)
		}
		return v
	case string:
		sanitized, valid := utils.SanitizeString(v)
		if !valid {
			return ""
		}
		return sanitized
	default:
		return v
	}
}

func processMultipartRequest(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		c.Abort()
		return
	}

	for key, values := range form.Value {
		for i, val := range values {
			sanitized, valid := utils.SanitizeString(val)
			if !valid {
				localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
				c.Abort()
				return
			}
			form.Value[key][i] = sanitized
		}
	}

	c.Request.MultipartForm = form
}
