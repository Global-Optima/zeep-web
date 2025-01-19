package data

import (
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/audit"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"net/http"
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

type DTO interface {
	any
}

type ExtendedDetails struct {
	BaseDetails
	DTO
}

func (d *ExtendedDetails) ToDetails() ([]byte, error) {
	return ToJSONB(d, true)
}

type MultipleItemDetails[T any] struct {
	IDs []uint `json:"ids"`
	DTO T
}

func (d *MultipleItemDetails[T]) ToDetails() ([]byte, error) {
	return ToJSONB(d, false)
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
		fields, err = audit.ExcludeEmptyJSONFields(input)
		if err != nil {
			return nil, err
		}
	} else {
		fields = input
	}
	logrus.Info(fields)

	return json.Marshal(fields)
}
