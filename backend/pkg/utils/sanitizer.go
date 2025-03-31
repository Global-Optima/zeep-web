package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/unicode/norm"
)

var (
	ErrReadBody            = errors.New("failed to read request body")
	ErrUnmarshalJSON       = errors.New("failed to unmarshal JSON")
	ErrParseMultipartForm  = errors.New("failed to parse multipart/form-data")
	ErrBindMultipartStruct = errors.New("failed to bind multipart form to struct")
	ErrSanitizeStruct      = errors.New("failed to sanitize struct data")
	ErrSoftSanitize        = errors.New("failed to soft-sanitize form field")
)

// ParseAndSanitize is the entry point that replicates the middleware logic
// in a handler-friendly function. It checks content-type, then either parses JSON
// or multipart data. If 'dto' is a pointer to a struct, it sanitizes that struct;
// otherwise, it falls back to a generic approach.
func ParseAndSanitize(c *gin.Context, dto interface{}) error {
	contentType := c.ContentType()
	switch {
	case strings.HasPrefix(contentType, "application/json"):
		return parseAndSanitizeJSON(c, dto)
	case strings.HasPrefix(contentType, "multipart/form-data"):
		return parseAndSanitizeMultipart(c, dto)
	default:
		// fallback or treat as JSON
		return parseAndSanitizeJSON(c, dto)
	}
}

// -----------------------------------------------------------
// JSON logic (mirrors processJSONRequest in your middleware)
// -----------------------------------------------------------

func parseAndSanitizeJSON(c *gin.Context, dto interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ErrReadBody
	}
	// Restore the original body so future reads still work
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Check if 'dto' is a pointer to a struct. If so, unmarshal into it.
	if isPointerToStruct(dto) {
		logrus.Info("Parsing JSON into provided DTO")
		// Unmarshal into the provided DTO
		if err := json.Unmarshal(body, dto); err != nil {
			return ErrUnmarshalJSON
		}
		// Sanitize
		if err := sanitizeStruct(dto); err != nil {
			return ErrSanitizeStruct
		}
		// Rewrite sanitized JSON back into request body
		sanitizedBody, _ := json.Marshal(dto)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
		return nil
	}

	// Fallback: parse into a generic interface{} and recursively sanitize
	var requestData interface{}
	if err := json.Unmarshal(body, &requestData); err != nil {
		return ErrUnmarshalJSON
	}
	requestData = sanitizeRecursive(requestData)
	sanitizedBody, _ := json.Marshal(requestData)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))

	return nil
}

// ----------------------------------------------------------------------
// Multipart logic (mirrors processMultipartRequest in your middleware)
// ----------------------------------------------------------------------

func parseAndSanitizeMultipart(c *gin.Context, dto interface{}) error {
	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(30 << 20); err != nil {
		return ErrParseMultipartForm
	}

	// If 'dto' is a pointer to struct, bind & sanitize
	if isPointerToStruct(dto) {
		if err := c.ShouldBind(dto); err != nil {
			// Could not bind to struct
			return ErrBindMultipartStruct
		}
		if err := sanitizeStruct(dto); err != nil {
			return ErrSanitizeStruct
		}
		// (Optional) update c.Request.MultipartForm with sanitized values
		updateMultipartFormFromDTO(c, dto)
		return nil
	}

	// Fallback: soft-sanitize all form fields
	form := c.Request.MultipartForm
	for key, values := range form.Value {
		for i, val := range values {
			sanitized, valid := SoftSanitizeString(val)
			if !valid {
				return ErrSoftSanitize
			}
			form.Value[key][i] = sanitized
		}
	}
	c.Request.MultipartForm = form
	return nil
}

// isPointerToStruct is a small helper that checks if 'obj' is a non-nil pointer to a struct
func isPointerToStruct(obj interface{}) bool {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr {
		return false
	}
	elem := rv.Elem()
	return elem.IsValid() && elem.Kind() == reflect.Struct
}

// -----------------------------------------------------------------------------
// Reflection-based sanitization (same as your middleware code).
// -----------------------------------------------------------------------------

func sanitizeStruct(dto interface{}) error {
	v := reflect.ValueOf(dto)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		// fallback if not pointer-to-struct
		logrus.Warn("sanitizeStruct called with non-pointer-to-struct")
		data, ok := dto.(*interface{})
		if ok && data != nil {
			*data = sanitizeRecursive(*data)
		}
		return nil
	}
	sanitizeRecursiveValue(v)
	return nil
}

// sanitizeRecursiveValue walks the struct fields looking for string/*string fields
func sanitizeRecursiveValue(v reflect.Value) {
	logrus.Debugf("Sanitizing value of kind: %s", v.Kind())
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			sanitizeRecursiveValue(v.Elem())
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)
			if !field.CanSet() {
				continue
			}
			bindingTag := fieldType.Tag.Get("binding")

			switch field.Kind() {
			case reflect.String:
				original := field.String()
				logrus.Info("Sanitizing string field:", fieldType.Name, "with binding tag:", bindingTag)
				field.SetString(sanitizeStringByBinding(original, bindingTag))

			case reflect.Ptr:
				// If it's a *string
				if field.Type().Elem().Kind() == reflect.String && !field.IsNil() {
					original := field.Elem().String()
					logrus.Info("Sanitizing *string field:", fieldType.Name, "with binding tag:", bindingTag)
					field.Elem().SetString(sanitizeStringByBinding(original, bindingTag))
				} else if !field.IsNil() {
					// pointer to something else - recurse
					logrus.Debugf("Recursing into pointer field: %s", fieldType.Name)
					sanitizeRecursiveValue(field.Elem())
				}
			case reflect.Slice, reflect.Array, reflect.Struct:
				sanitizeRecursiveValue(field)
			default:
				// skip numeric, bool, etc.
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			sanitizeRecursiveValue(v.Index(i))
		}
	}
}

// sanitizeStringByBinding applies strict or soft sanitize based on the "binding" tag
func sanitizeStringByBinding(original, bindingTag string) string {
	if original == "" {
		return original
	}
	if strings.Contains(bindingTag, "min=") {
		logrus.Info("Applying strict sanitization for min binding")
		sanitized, valid := SanitizeString(original)
		if !valid {
			logrus.Warnf("Sanitization failed for string: %s", original)
			return ""
		}
		logrus.Debugf("Sanitized string: %s", sanitized)
		return sanitized
	}
	if strings.Contains(bindingTag, "omitempty") {
		logrus.Info("Applying soft sanitization for omitempty binding")
		if original == "" {
			return original
		}
		soft, _ := SoftSanitizeString(original)
		logrus.Debugf("Soft sanitized string: %s", soft)
		return soft
	}
	// Default: strict
	sanitized, valid := SanitizeString(original)
	if !valid {
		return ""
	}
	return sanitized
}

// sanitizeRecursive is the fallback for raw data that isn't part of your struct
func sanitizeRecursive(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			v[key] = sanitizeRecursive(val)
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
		sanitized, valid := SanitizeString(v)
		if !valid {
			return ""
		}
		return sanitized
	default:
		return v
	}
}

// updateMultipartFormFromDTO rewrites the sanitized string fields in c.Request.MultipartForm
func updateMultipartFormFromDTO(c *gin.Context, dto interface{}) {
	if c.Request.MultipartForm == nil {
		return
	}
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

// ----------------------------------------------------------------------------
// The rest is your existing string-sanitizing and validator logic
// ----------------------------------------------------------------------------

func SanitizeString(input string) (string, bool) {
	if !utf8.ValidString(input) {
		return "", false
	}
	normalizedUTF8 := norm.NFC.String(input)
	trimmed := strings.TrimSpace(normalizedUTF8)
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized := spaceRegex.ReplaceAllString(trimmed, " ")
	normalized = removeInvisibleCharacters(normalized)

	// If there's nothing left, we consider it invalid
	if len(normalized) == 0 {
		return "", false
	}
	return normalized, true
}

func SoftSanitizeString(input string) (string, bool) {
	if !utf8.ValidString(input) {
		return "", false
	}
	normalizedUTF8 := norm.NFC.String(input)
	trimmed := strings.TrimSpace(normalizedUTF8)
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized := spaceRegex.ReplaceAllString(trimmed, " ")
	normalized = removeInvisibleCharacters(normalized)
	return normalized, true
}

func removeInvisibleCharacters(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if unicode.IsPrint(char) && !unicode.IsControl(char) {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

// Example validator usage
func ValidateSanitizedString(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	_, valid := SanitizeString(value)
	return valid
}

func RegisterCustomValidators(validate *validator.Validate) {
	err := validate.RegisterValidation("customSanitize", ValidateSanitizedString)
	if err != nil {
		fmt.Printf("failed to register custom validator: %v", err)
		return
	}
}

func InitValidators() {
	validate := validator.New()
	RegisterCustomValidators(validate)
}

func RoundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

func MergeDistinct[T comparable](arr1, arr2 []T) []T {
	uniqueMap := make(map[T]struct{})

	for _, v := range arr1 {
		uniqueMap[v] = struct{}{}
	}
	for _, v := range arr2 {
		uniqueMap[v] = struct{}{}
	}

	mergedSlice := make([]T, 0, len(uniqueMap))
	for key := range uniqueMap {
		mergedSlice = append(mergedSlice, key)
	}

	return mergedSlice
}
