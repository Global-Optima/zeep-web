package data

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"net/http"
	"reflect"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type OperationType string

const (
	CreateMultipleOperation OperationType = "CREATE MULTIPLE"
	CreateOperation         OperationType = "CREATE"
	UpdateOperation         OperationType = "UPDATE"
	DeleteOperation         OperationType = "DELETE"
)

type ComponentName string

const (
	ProductComponent             ComponentName = "PRODUCT"
	StoreProductComponent        ComponentName = "STORE PRODUCT"
	EmployeeComponent            ComponentName = "EMPLOYEE"
	AdditiveComponent            ComponentName = "ADDITIVE"
	StoreAdditiveComponent       ComponentName = "STORE ADDITIVE"
	ProductSizeComponent         ComponentName = "PRODUCT SIZE"
	StoreProductSizeComponent    ComponentName = "STORE PRODUCT SIZE"
	RecipeStepsComponent         ComponentName = "RECIPE STEPS"
	StoreComponent               ComponentName = "STORE"
	WarehouseComponent           ComponentName = "WAREHOUSE"
	StoreWarehouseStockComponent ComponentName = "STORE WAREHOUSE STOCK"
	IngredientComponent          ComponentName = "INGREDIENT"
	StoreProductSizesComponent   ComponentName = "STORE PRODUCT SIZES"
)

func (o OperationType) ToString() string {
	return string(o)
}

func (c ComponentName) ToString() string {
	return string(c)
}

type AuditDetails interface {
	ToDetails() ([]byte, error)
}

type BaseDetails struct {
	ID uint `json:"id"`
}

func (b *BaseDetails) ToDetails() ([]byte, error) {
	return ToJSONB(b, false)
}

type ItemDetails[T CustomFields] struct {
	BaseDetails
	CustomFields T
}

type MultipleCreationDetails[T CustomFields] struct {
	BaseDetails  []BaseDetails
	CustomFields T
}

func (d *MultipleCreationDetails[T]) ToDetails() ([]byte, error) {
	return ToJSONB(d, false)
}

type CustomFields interface {
	any
}

type UpdateDetails[T CustomFields] struct {
	ItemDetails[T]
}

func (d *UpdateDetails[T]) ToDetails() ([]byte, error) {
	return ToJSONB(d, true)
}

type AuditAction struct {
	OperationType OperationType
	ComponentName ComponentName
}

func (a AuditAction) ToString() string {
	return a.OperationType.ToString() + " " + a.ComponentName.ToString()
}

type HTTPMethod string

func (m HTTPMethod) ToString() string {
	return string(m)
}

func ToHTTPMethod(method string) (HTTPMethod, error) {
	httpMethod := HTTPMethod(method)
	switch httpMethod {
	case http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodPost:
		return httpMethod, nil
	}
	return "", errors.New("invalid http method")
}

type FieldChange struct {
	FieldName string `json:"fieldName"`
	Value     string `json:"value"`
}

type EmployeeAudit struct {
	BaseEntity
	EmployeeID    uint           `gorm:"index;not null"`
	Employee      Employee       `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employees"`
	OperationType OperationType  `gorm:"type:operation_type;not null" sort:"operationType"`
	ComponentName ComponentName  `gorm:"type:component_name;not null" sort:"componentName"`
	Details       datatypes.JSON `gorm:"type:jsonb"`
	IPAddress     string         `gorm:"column:ip_address;type:varchar(45);not null"`
	ResourceUrl   string         `gorm:"column:resource_url;type:text;not null"`
	Method        HTTPMethod     `gorm:"type:http_method;not null" sort:"method"`
}

func ToJSONB(input interface{}, excludeEmptyFields bool) ([]byte, error) {
	var fields interface{}
	var err error

	if excludeEmptyFields {
		fields, err = ExcludeEmptyJSONFields(input)
		if err != nil {
			return nil, err
		}
	} else {
		fields = input
	}
	logrus.Info(fields)

	return json.Marshal(fields)
}

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
		if fieldType.Anonymous || jsonName == "CustomFields" {
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
