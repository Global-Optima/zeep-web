package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const DTO_KEY = "dto"

func WithDTO(dto interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DTO_KEY, dto)
		c.Next()
	}
}

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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		c.Abort()
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if dto, exists := c.Get(DTO_KEY); exists {
		if err := json.Unmarshal(body, dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}

		if err := sanitizeStruct(dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}

		sanitizedBody, _ := json.Marshal(dto)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
	} else {
		logrus.Infoln("Sanitizing fallback: request data")
		var requestData interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}
		requestData = sanitizeRecursive(requestData)
		sanitizedBody, _ := json.Marshal(requestData)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
	}
}

func sanitizeStruct(dto interface{}) error {
	v := reflect.ValueOf(dto)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("sanitize")
		if tag == "skip" {
			logrus.Infoln("Skipping sanitization for field", fieldType.Name)
			continue
		}

		// Use soft sanitization if tag is "soft"
		if tag == "soft" {
			logrus.Infoln("Soft sanitizing field", fieldType.Name)
			if field.Kind() == reflect.String && field.CanSet() {
				original := field.String()
				softSanitized, _ := utils.SoftSanitizeString(original)
				field.SetString(softSanitized)
				continue
			}

			if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.String && !field.IsNil() {
				original := field.Elem().String()
				softSanitized, _ := utils.SoftSanitizeString(original)
				field.Elem().SetString(softSanitized)
				continue
			}
		}

		// Default sanitization for other fields
		if field.Kind() == reflect.String && field.CanSet() {
			original := field.String()
			sanitized, valid := utils.SanitizeString(original)
			if !valid {
				field.SetString("")
			} else {
				field.SetString(sanitized)
			}
		}

		if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.String && !field.IsNil() {
			original := field.Elem().String()
			sanitized, valid := utils.SanitizeString(original)
			if !valid {
				field.Elem().SetString("")
			} else {
				field.Elem().SetString(sanitized)
			}
		}
	}
	return nil
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

func softSanitizeRecursive(data interface{}) interface{} {
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
		sanitized, valid := utils.SoftSanitizeString(v)
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
