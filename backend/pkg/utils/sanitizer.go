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
	ErrReadBody             = errors.New("failed to read request body")
	ErrUnmarshalJSON        = errors.New("failed to unmarshal JSON")
	ErrParseMultipartForm   = errors.New("failed to parse multipart/form-data")
	ErrBindMultipartStruct  = errors.New("failed to bind multipart form to struct")
	ErrSanitizeStruct       = errors.New("failed to sanitize struct data")
	ErrSoftSanitize         = errors.New("failed to soft-sanitize form field")
	ErrStrictSanitizeFailed = errors.New("strict sanitization failed (invalid or empty)")
)

func ParseRequestBody(c *gin.Context, dto interface{}) error {
	contentType := c.ContentType()
	switch {
	case strings.HasPrefix(contentType, "application/json"):
		return ParseRequestBodyJSON(c, dto)
	case strings.HasPrefix(contentType, "multipart/form-data"):
		return ParseRequestBodyMultipart(c, dto)
	default:
		// fallback treat as JSON
		return ParseRequestBodyJSON(c, dto)
	}
}

func ParseRequestBodyJSON(c *gin.Context, dto interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ErrReadBody
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if isPointerToStruct(dto) {
		if err := json.Unmarshal(body, dto); err != nil {
			logrus.Errorf("error unmarshalling JSON: %v", err)
			return ErrUnmarshalJSON
		}

		// Sanitize
		if err := sanitizeStruct(dto); err != nil {
			return fmt.Errorf("%w: %v", ErrSanitizeStruct, err)
		}

		sanitizedBody, _ := json.Marshal(dto)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
		return nil
	}

	var requestData interface{}
	if err := json.Unmarshal(body, &requestData); err != nil {
		return ErrUnmarshalJSON
	}
	sanitized := sanitizeRecursive(requestData)
	sanitizedBody, _ := json.Marshal(sanitized)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))

	return nil
}

func ParseRequestBodyMultipart(c *gin.Context, dto interface{}) error {
	if err := c.Request.ParseMultipartForm(30 << 20); err != nil {
		return ErrParseMultipartForm
	}

	if isPointerToStruct(dto) {
		if err := c.ShouldBind(dto); err != nil {
			return ErrBindMultipartStruct
		}
		if err := sanitizeStruct(dto); err != nil {
			return fmt.Errorf("%w: %v", ErrSanitizeStruct, err)
		}

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

func sanitizeStruct(dto interface{}) error {
	v := reflect.ValueOf(dto)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		// fallback if not pointer-to-struct
		data, ok := dto.(*interface{})
		if ok && data != nil {
			*data = sanitizeRecursive(*data)
		}
		return nil
	}

	if err := sanitizeRecursiveValue(v); err != nil {
		return err
	}
	return nil
}

func sanitizeRecursiveValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			return sanitizeRecursiveValue(v.Elem())
		}
		return nil

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
				sanitized, err := sanitizeStringByBinding(original, bindingTag)
				if err != nil {
					return err // bubble up
				}
				field.SetString(sanitized)

			case reflect.Ptr:
				// If it's a *string
				if field.Type().Elem().Kind() == reflect.String {
					if field.IsNil() {
						empty := ""
						field.Set(reflect.ValueOf(&empty))
					} else {
						original := field.Elem().String()
						sanitized, err := sanitizeStringByBinding(original, bindingTag)
						if err != nil {
							return err
						}
						field.Elem().SetString(sanitized)
					}
				} else if !field.IsNil() {
					// pointer to something else - recurse
					if err := sanitizeRecursiveValue(field.Elem()); err != nil {
						return err
					}
				}

			case reflect.Slice, reflect.Array, reflect.Struct:
				if err := sanitizeRecursiveValue(field); err != nil {
					return err
				}
			default:
				// skip numeric, bool, etc.
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := sanitizeRecursiveValue(v.Index(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

// sanitizeStringByBinding returns (sanitizedString, error).
// If "strict" fails, we return ErrStrictSanitizeFailed.
func sanitizeStringByBinding(original, bindingTag string) (string, error) {
	if original == "" {
		return original, nil
	}

	// If the binding tag has "min=", treat as strict
	if strings.Contains(bindingTag, "min=") {
		sanitized, valid := SanitizeString(original)
		if !valid {
			return "", ErrStrictSanitizeFailed
		}
		return sanitized, nil
	}

	// If the binding tag has "omitempty" (but no "min="), do soft
	if strings.Contains(bindingTag, "omitempty") {
		if original == "" {
			return original, nil
		}
		soft, _ := SoftSanitizeString(original)
		return soft, nil
	}

	// Default: strict
	sanitized, valid := SanitizeString(original)
	if !valid {
		return "", ErrStrictSanitizeFailed
	}
	return sanitized, nil
}

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

func SanitizeString(input string) (string, bool) {
	if !utf8.ValidString(input) {
		return "", false
	}
	normalizedUTF8 := norm.NFC.String(input)
	trimmed := strings.TrimSpace(normalizedUTF8)
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized := spaceRegex.ReplaceAllString(trimmed, " ")
	normalized = removeInvisibleCharacters(normalized)

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
	_ = validate.RegisterValidation("customSanitize", ValidateSanitizedString)
}

func InitValidators() {
	validate := validator.New()
	RegisterCustomValidators(validate)
}

func RoundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

func UnionSlices[T comparable](arr1, arr2 []T) []T {
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

func DiffSlice[T comparable](all, subset []T) []T {
	m := make(map[T]struct{}, len(subset))
	for _, v := range subset {
		m[v] = struct{}{}
	}

	var diff []T
	for _, v := range all {
		if _, exists := m[v]; !exists {
			diff = append(diff, v)
		}
	}
	return diff
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
