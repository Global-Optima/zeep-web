package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"

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
		if strings.HasPrefix(contentType, "application/json") {
			processJSONRequest(c)
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
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
	if !strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		return
	}

	if err := c.Request.ParseMultipartForm(30 << 20); err != nil {
		localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
		c.Abort()
		return
	}

	if dto, exists := c.Get(DTO_KEY); exists {
		if err := c.ShouldBind(dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}
		if err := sanitizeStruct(dto); err != nil {
			localization.SendLocalizedResponseWithKey(c, localization.ErrMessageBindingJSON)
			c.Abort()
			return
		}
		updateMultipartFormFromDTO(c, dto)
		return
	}

	// Fallback: without DTO we can only apply soft sanitation.
	form := c.Request.MultipartForm
	for key, values := range form.Value {
		for i, val := range values {
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

// updateMultipartFormFromDTO updates c.Request.MultipartForm.Value based on the sanitized DTO.
// It uses reflection to iterate over the DTO's fields, uses the "form" tag to look up
// the corresponding key in the multipart form, and writes the sanitized value.
func updateMultipartFormFromDTO(c *gin.Context, dto interface{}) {
	form := c.Request.MultipartForm
	v := reflect.ValueOf(dto).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if !field.CanInterface() {
			continue
		}
		formTag := fieldType.Tag.Get("form")
		if formTag == "" {
			continue
		}

		var val string
		switch field.Kind() {
		case reflect.String:
			val = field.String()
		case reflect.Ptr:
			if field.Type().Elem().Kind() == reflect.String && !field.IsNil() {
				val = field.Elem().String()
			} else {
				continue
			}
		default:
			continue
		}
		form.Value[formTag] = []string{val}
	}
	c.Request.MultipartForm = form
}

// sanitizeStruct recursively processes a DTO struct (or nested structs).
// It expects dto to be a pointer to a struct.
func sanitizeStruct(dto interface{}) error {
	v := reflect.ValueOf(dto)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil
	}
	sanitizeRecursiveValue(v)
	return nil
}

// sanitizeRecursiveValue walks through any value recursively and sanitizes string fields.
// It checks the "binding" tag:
// if the tag contains "min=", strict sanitation is applied;
// if the tag contains "omitempty" (and not "min="), soft sanitation is applied;
// otherwise, strict sanitation is used.
func sanitizeRecursiveValue(v reflect.Value) {
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
			if !field.CanSet() {
				continue
			}
			bindingTag := fieldType.Tag.Get("binding")
			if field.Kind() == reflect.String {
				original := field.String()
				field.SetString(sanitizeStringByBinding(original, bindingTag))
				continue
			}
			if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.String && !field.IsNil() {
				original := field.Elem().String()
				field.Elem().SetString(sanitizeStringByBinding(original, bindingTag))
				continue
			}
			sanitizeRecursiveValue(field)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			sanitizeRecursiveValue(v.Index(i))
		}
	}
}

// sanitizeStringByBinding applies the appropriate sanitation based on the binding tag.
// If the binding tag contains "min=", strict sanitation is applied (empty strings are disallowed after trim).
// If it contains "omitempty" (and no "min="), soft sanitation is applied (empty strings remain as-is).
// Otherwise, strict sanitation is applied.
func sanitizeStringByBinding(original, bindingTag string) string {
	if original == "" {
		return original
	}
	if strings.Contains(bindingTag, "min=") {
		sanitized, valid := utils.SanitizeString(original)
		if !valid {
			return ""
		}
		return sanitized
	}
	if strings.Contains(bindingTag, "omitempty") {
		if original == "" {
			return original
		}
		soft, _ := utils.SoftSanitizeString(original)
		return soft
	}
	sanitized, valid := utils.SanitizeString(original)
	if !valid {
		return ""
	}
	return sanitized
}

// sanitizeRecursive is the fallback for non-DTO JSON data, applying strict sanitation on strings.
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
