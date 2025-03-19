package audit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func ExcludeEmptyJSONFields(input interface{}) (map[string]interface{}, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	val := reflect.ValueOf(input)

	// Unwrap multiple pointer levels safely
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, errors.New("nil pointer encountered while unwrapping")
		}
		val = val.Elem()
	}

	// Ensure final value is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct or pointer-to-struct, got %s", val.Kind())
	}

	typ := val.Type()
	changes := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		// Parse JSON tag
		jsonTag := fieldType.Tag.Get("json")
		jsonName := strings.Split(jsonTag, ",")[0]
		// Use the field name if no JSON name is provided
		if jsonName == "" || jsonName == "-" {
			jsonName = fieldType.Name
		}

		// Handle embedded or anonymous fields and CustomFields
		if fieldType.Anonymous || jsonName == "DTO" {
			embeddedFields, err := ExcludeEmptyJSONFields(fieldValue.Interface())
			if err != nil {
				continue
			}

			// Merge embedded fields directly into the parent map
			for key, value := range embeddedFields {
				changes[key] = value
			}
			continue
		}

		// Handle nested structs or pointers to structs
		if fieldValue.Kind() == reflect.Struct ||
			(fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {

			nested, err := ExcludeEmptyJSONFields(fieldValue.Interface())
			if err != nil {
				continue
			}

			// Only add this field if the nested map is not empty
			if len(nested) > 0 {
				changes[jsonName] = nested
			}
			continue
		}

		// Non-struct field: do zero-check
		zero := reflect.Zero(fieldValue.Type()).Interface()
		current := fieldValue.Interface()

		if !reflect.DeepEqual(current, zero) {
			changes[jsonName] = current
		}
	}

	return changes, nil
}
