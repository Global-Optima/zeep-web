package data

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils/audit"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type OperationType string

const (
	GetOperation    OperationType = "GET"
	CreateOperation OperationType = "CREATE"
	UpdateOperation OperationType = "UPDATE"
	DeleteOperation OperationType = "DELETE"
)

type ComponentName string

const (
	FranchiseeComponent            ComponentName = "FRANCHISEE"
	RegionComponent                ComponentName = "REGION"
	ProductComponent               ComponentName = "PRODUCT"
	ProductCategoryComponent       ComponentName = "PRODUCT_CATEGORY"
	StoreProductComponent          ComponentName = "STORE_PRODUCT"
	EmployeeComponent              ComponentName = "EMPLOYEE"
	StoreEmployeeComponent         ComponentName = "STORE_EMPLOYEE"
	WarehouseEmployeeComponent     ComponentName = "WAREHOUSE_EMPLOYEE"
	FranchiseeEmployeeComponent    ComponentName = "FRANCHISEE_EMPLOYEE"
	RegionEmployeeComponent        ComponentName = "REGION_EMPLOYEE"
	AdminEmployeeComponent         ComponentName = "ADMIN_EMPLOYEE"
	AdditiveComponent              ComponentName = "ADDITIVE"
	AdditiveCategoryComponent      ComponentName = "ADDITIVE_CATEGORY"
	StoreAdditiveComponent         ComponentName = "STORE_ADDITIVE"
	ProductSizeComponent           ComponentName = "PRODUCT_SIZE"
	RecipeStepsComponent           ComponentName = "RECIPE_STEPS"
	StoreComponent                 ComponentName = "STORE"
	WarehouseComponent             ComponentName = "WAREHOUSE"
	StoreStockComponent            ComponentName = "STORE_STOCK"
	IngredientComponent            ComponentName = "INGREDIENT"
	IngredientCategoryComponent    ComponentName = "INGREDIENT_CATEGORY"
	StockRequestComponent          ComponentName = "STOCK_REQUESTS"
	StockMaterialComponent         ComponentName = "STOCK_MATERIAL"
	StockMaterialCategoryComponent ComponentName = "STOCK_MATERIAL_CATEGORY"
	WarehouseStockComponent        ComponentName = "WAREHOUSE_STOCK"
	SupplierComponent              ComponentName = "SUPPLIER"
	UnitComponent                  ComponentName = "UNIT"
	OrderComponent                 ComponentName = "ORDER"

	AuthenticationComponent ComponentName = "AUTH"
)

func (o OperationType) ToString() string {
	return string(o)
}

func (c ComponentName) ToString() string {
	return string(c)
}

type AuditDetails interface {
	ToDetails() ([]byte, error)
	GetBaseDetails() *BaseDetails
}

type BaseDetails struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (b *BaseDetails) ToDetails() ([]byte, error) {
	return ToJSONB(b, false)
}

func (b *BaseDetails) GetBaseDetails() *BaseDetails {
	return b
}

type DTO interface {
	any
}

type ExtendedDetails struct {
	BaseDetails
	DTO `json:"data"`
}

func (d *ExtendedDetails) GetBaseDetails() *BaseDetails {
	return &d.BaseDetails
}

func (d *ExtendedDetails) ToDetails() ([]byte, error) {
	return ToJSONB(d, true)
}

type StoreInfo struct {
	StoreID   uint   `json:"storeId"`
	StoreName string `json:"storeName"`
}

type WarehouseInfo struct {
	WarehouseID   uint   `json:"warehouseId"`
	WarehouseName string `json:"warehouseName"`
}

type FranchiseeInfo struct {
	FranchiseeID   uint   `json:"franchiseeId"`
	FranchiseeName string `json:"franchiseeName"`
}

type RegionInfo struct {
	RegionID   uint   `json:"regionId"`
	RegionName string `json:"regionName"`
}

type ExtendedDetailsStore struct {
	ExtendedDetails
	StoreInfo
}

func (d *ExtendedDetailsStore) ToDetails() ([]byte, error) {
	return ToJSONB(d, false)
}

func (d *ExtendedDetailsStore) SetStoreName(name string) {
	d.StoreInfo.StoreName = name
}

type ExtendedDetailsWarehouse struct {
	ExtendedDetails
	WarehouseInfo
}

func (d *ExtendedDetailsWarehouse) ToDetails() ([]byte, error) {
	return ToJSONB(d, false)
}

func (d *ExtendedDetailsWarehouse) SetWarehouseName(name string) {
	d.WarehouseInfo.WarehouseName = name
}

type ExtendedDetailsFranchisee struct {
	ExtendedDetails
	FranchiseeInfo
}

func (d *ExtendedDetailsFranchisee) ToDetails() ([]byte, error) {
	return ToJSONB(d, false)
}

func (d *ExtendedDetailsFranchisee) SetFranchiseeName(name string) {
	d.FranchiseeInfo.FranchiseeName = name
}

type ExtendedDetailsRegion struct {
	ExtendedDetails
	RegionInfo
}

func (d *ExtendedDetailsRegion) ToDetails() ([]byte, error) {
	return ToJSONB(d, false)
}

func (d *ExtendedDetailsRegion) SetRegionName(name string) {
	d.RegionInfo.RegionName = name
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

	return json.Marshal(fields)
}
