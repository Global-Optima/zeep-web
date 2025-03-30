package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

const DTO_KEY = "dto"

// WithDTO sets the DTO in the context for later sanitization.
func WithDTO(dto interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DTO_KEY, dto)
		c.Next()
	}
}

// SanitizeMiddleware handles sanitization for JSON and multipart/form-data requests.
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
	// Restore the original body for further use.
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if dto, exists := c.Get(DTO_KEY); exists {
		if err := json.Unmarshal(body, dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}

		// Recursively sanitize the DTO.
		if err := sanitizeStruct(dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}

		sanitizedBody, _ := json.Marshal(dto)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
	} else {
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

func processMultipartRequest(c *gin.Context) {
	// Parse multipart form (adjust size as needed)
	if err := c.Request.ParseMultipartForm(30 << 20); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		c.Abort()
		return
	}

	// If a DTO is present, bind form values into it and sanitize.
	if dto, exists := c.Get(DTO_KEY); exists {
		// Bind multipart form values into the DTO.
		if err := c.ShouldBind(dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}
		// Recursively sanitize the DTO.
		if err := sanitizeStruct(dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}
		// DTO-bound data is now sanitized—exit processing.
		return
	}

	// Fallback: manually iterate over form values.
	form := c.Request.MultipartForm
	for key, values := range form.Value {
		for i, val := range values {
			// Use soft sanitation for all fields; empty values will be preserved.
			sanitized, valid := utils.SoftSanitizeString(val)
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

// sanitizeStruct recursively processes a DTO struct (or nested structs) based on tags.
// Expects dto to be a pointer to a struct.
func sanitizeStruct(dto interface{}) error {
	v := reflect.ValueOf(dto)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil
	}
	sanitizeRecursiveValue(v)
	return nil
}

// sanitizeRecursiveValue walks any value recursively and sanitizes string fields.
// It respects the "sanitize" tag if provided ("skip" to leave unchanged).
// Otherwise, it applies full sanitization.
// If the string is empty, it is returned as-is (allowing empty values via "omitempty").
func sanitizeRecursiveValue(v reflect.Value) {
	// If pointer, process its element (if non-nil).
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		sanitizeRecursiveValue(v.Elem())
		return
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)
			tag := fieldType.Tag.Get("sanitize")
			if tag == "skip" {
				continue
			}

			// If the field is a string.
			if field.Kind() == reflect.String && field.CanSet() {
				original := field.String()
				field.SetString(sanitizeStringByTag(original, tag))
				continue
			}

			// If the field is a pointer to a string.
			if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.String && !field.IsNil() {
				original := field.Elem().String()
				field.Elem().SetString(sanitizeStringByTag(original, tag))
				continue
			}

			// Recursively process nested structs, pointers, slices, etc.
			sanitizeRecursiveValue(field)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			sanitizeRecursiveValue(v.Index(i))
		}
	}
}

// sanitizeStringByTag applies the sanitization method.
// If the original string is empty, it returns it as-is.
func sanitizeStringByTag(original, tag string) string {
	if original == "" {
		return original
	}
	// If you need to support additional tag logic (e.g. "soft"),
	// you can add conditions here.
	sanitized, valid := utils.SanitizeString(original)
	if !valid {
		return ""
	}
	return sanitized
}

// sanitizeRecursive is the fallback for non-DTO JSON data, applying default sanitation on strings.
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
		// If the string is empty, return it as-is.
		if v == "" {
			return v
		}
		sanitized, valid := utils.SanitizeString(v)
		if !valid {
			return ""
		}
		return sanitized
	default:
		return v
	}
}
