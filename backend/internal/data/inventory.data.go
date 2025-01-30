package data

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StockRequestStatus string

var (
	StockRequestCreated             StockRequestStatus = "CREATED"
	StockRequestProcessed           StockRequestStatus = "PROCESSED"
	StockRequestInDelivery          StockRequestStatus = "IN_DELIVERY"
	StockRequestCompleted           StockRequestStatus = "COMPLETED"
	StockRequestRejectedByStore     StockRequestStatus = "REJECTED_BY_STORE"
	StockRequestRejectedByWarehouse StockRequestStatus = "REJECTED_BY_WAREHOUSE"
	StockRequestAcceptedWithChange  StockRequestStatus = "ACCEPTED_WITH_CHANGE"
)

type Region struct {
	BaseEntity
	Name           string           `gorm:"size:255;not null" sort:"name"`
	Warehouses     []Warehouse      `gorm:"foreignKey:RegionID;constraint:OnDelete:CASCADE"`
	RegionManagers []RegionEmployee `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

type Warehouse struct {
	BaseEntity
	FacilityAddressID uint            `gorm:"not null;index"`
	RegionID          uint            `gorm:"not null;index"`
	FacilityAddress   FacilityAddress `gorm:"foreignKey:FacilityAddressID;constraint:OnDelete:CASCADE"`
	Name              string          `gorm:"size:255;not null" sort:"name"`
}

type StoreWarehouse struct {
	BaseEntity
	StoreID     uint      `gorm:"not null;index"`
	Store       Store     `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE" sort:"stores"`
	WarehouseID uint      `gorm:"not null;index"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE" sort:"warehouses"`
}

type StoreWarehouseStock struct {
	BaseEntity
	StoreWarehouseID  uint           `gorm:"not null;index"`
	StoreWarehouse    StoreWarehouse `gorm:"foreignKey:StoreWarehouseID;constraint:OnDelete:CASCADE"`
	IngredientID      uint           `gorm:"not null;index"`
	Ingredient        Ingredient     `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE" sort:"ingredients"`
	LowStockThreshold float64        `gorm:"type:decimal(10,2);not null;check:low_stock_threshold > 0" sort:"lowStockThreshold"`
	Quantity          float64        `gorm:"type:decimal(10,2);not null;check:quantity >= 0" sort:"quantity"`
}

type StockRequest struct {
	BaseEntity
	StoreID          uint                     `gorm:"not null;index"` // Links to Store
	Store            Store                    `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	WarehouseID      uint                     `gorm:"not null;index"` // Central warehouse fulfilling the request
	Warehouse        Warehouse                `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE" sort:"warehouses"`
	Status           StockRequestStatus       `gorm:"size:50;not null" sort:"status"`
	Details          []datatypes.JSON         `gorm:"type:jsonb"`
	StoreComment     *string                  `gorm:"type:text"` // Store-specific comments
	WarehouseComment *string                  `gorm:"type:text"` // Warehouse-specific comments
	Ingredients      []StockRequestIngredient `gorm:"foreignKey:StockRequestID;constraint:OnDelete:CASCADE"`
}

type StockRequestIngredient struct {
	BaseEntity
	StockRequestID  uint          `gorm:"not null;index"`
	StockRequest    StockRequest  `gorm:"foreignKey:StockRequestID;constraint:OnDelete:CASCADE"`
	IngredientID    uint          `gorm:"not null;index"`
	Ingredient      Ingredient    `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	StockMaterialID uint          `gorm:"not null;index"` // Selected stock material
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	Quantity        float64       `gorm:"type:decimal(10,2);not null;check:quantity > 0" sort:"quantity"`
	DeliveredDate   time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP" sort:"deliveryDate"` // Delivery start date
	ExpirationDate  time.Time     `gorm:"not null" sort:"expirationDate"`                         // Calculated from DeliveredDate + ExpirationPeriodInDays
}

// Hooks for StockRequestIngredient
func (s *StockRequestIngredient) BeforeCreate(tx *gorm.DB) (err error) {
	s.DeliveredDate = toUTC(s.DeliveredDate)
	s.ExpirationDate = toUTC(s.ExpirationDate)
	return
}

func (s *StockRequestIngredient) BeforeUpdate(tx *gorm.DB) (err error) {
	s.DeliveredDate = toUTC(s.DeliveredDate)
	s.ExpirationDate = toUTC(s.ExpirationDate)
	return
}

func (s *StockRequestIngredient) AfterFind(tx *gorm.DB) (err error) {
	s.DeliveredDate = toUTC(s.DeliveredDate)
	s.ExpirationDate = toUTC(s.ExpirationDate)
	return
}

type Supplier struct {
	BaseEntity
	Name         string `gorm:"size:255;not null" sort:"name"`
	ContactEmail string `gorm:"size:255"`
	ContactPhone string `gorm:"size:16;uniqueIndex"`
	City         string `gorm:"size:100;not null"`
	Address      string `gorm:"size:255"`
}

type WarehouseStock struct {
	BaseEntity
	WarehouseID     uint          `gorm:"not null;index"`
	Warehouse       Warehouse     `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE" sort:"warehouses"`
	StockMaterialID uint          `gorm:"index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE" sort:"stockMaterials"`
	Quantity        float64       `gorm:"type:decimal(10,2);not null;check:quantity >= 0" sort:"name"`
}

type StockMaterial struct {
	BaseEntity
	Name                   string                `gorm:"size:255;not null" sort:"name"`
	Description            string                `gorm:"type:text"`
	IngredientID           uint                  `gorm:"not null;index"` // Link to the ingredient
	Ingredient             Ingredient            `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	SafetyStock            float64               `gorm:"type:decimal(10,2);not null" sort:"safetyStock"`
	UnitID                 uint                  `gorm:"not null"`
	Unit                   Unit                  `gorm:"foreignKey:UnitID;constraint:OnDelete:SET NULL"`
	Size                   float64               `gorm:"type:decimal(10,2);not null"`
	CategoryID             uint                  `gorm:"not null"` // Link to IngredientCategory
	StockMaterialCategory  StockMaterialCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	Barcode                string                `gorm:"unique;size:255"`
	ExpirationPeriodInDays int                   `gorm:"not null;default:1095" sort:"expirationPeriodInDays"` // 3 years in days
	IsActive               bool                  `gorm:"not null;default:true" sort:"isActive"`
}

type StockMaterialCategory struct {
	BaseEntity
	Name           string          `gorm:"size:255;not null;uniqueIndex"`
	Description    string          `gorm:"type:text"`
	StockMaterials []StockMaterial `gorm:"foreignKey:CategoryID"`
}

type Unit struct {
	BaseEntity
	Name             string  `gorm:"size:50;not null"`
	ConversionFactor float64 `gorm:"type:decimal(10,4);not null"` // To base unit
}

type SupplierMaterial struct {
	BaseEntity
	StockMaterialID uint            `gorm:"not null;index"`
	StockMaterial   StockMaterial   `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	SupplierID      uint            `gorm:"not null;index"`
	Supplier        Supplier        `gorm:"foreignKey:SupplierID;constraint:OnDelete:CASCADE"`
	SupplierPrices  []SupplierPrice `gorm:"foreignKey:SupplierMaterialID;constraint:OnDelete:CASCADE"`
}

type SupplierPrice struct {
	BaseEntity
	SupplierMaterialID uint             `gorm:"not null;index"`
	SupplierMaterial   SupplierMaterial `gorm:"foreignKey:SupplierMaterialID;constraint:OnDelete:CASCADE"`
	BasePrice          float64          `gorm:"type:decimal(10,2);not null"`
}

type SupplierWarehouseDelivery struct {
	BaseEntity
	SupplierID   uint                                `gorm:"not null"`
	Supplier     Supplier                            `gorm:"foreignKey:SupplierID;constraint:OnDelete:CASCADE"`
	WarehouseID  uint                                `gorm:"not null"`
	Warehouse    Warehouse                           `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE"`
	Materials    []SupplierWarehouseDeliveryMaterial `gorm:"foreignKey:DeliveryID;constraint:OnDelete:CASCADE"`
	DeliveryDate time.Time                           `gorm:"not null;default:CURRENT_TIMESTAMP" sort:"deliveryDate"`
}

// Hooks for SupplierWarehouseDelivery
func (s *SupplierWarehouseDelivery) BeforeCreate(tx *gorm.DB) (err error) {
	s.DeliveryDate = toUTC(s.DeliveryDate)
	return
}

func (s *SupplierWarehouseDelivery) AfterFind(tx *gorm.DB) (err error) {
	s.DeliveryDate = toUTC(s.DeliveryDate)
	return
}

type SupplierWarehouseDeliveryMaterial struct {
	BaseEntity
	DeliveryID      uint                      `gorm:"not null;index"`
	Delivery        SupplierWarehouseDelivery `gorm:"foreignKey:DeliveryID;constraint:OnDelete:CASCADE"`
	StockMaterialID uint                      `gorm:"not null;index"`
	StockMaterial   StockMaterial             `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	Barcode         string                    `gorm:"size:255;not null"`
	Quantity        float64                   `gorm:"type:decimal(10,2);not null;check:quantity > 0" sort:"quantity"`
	ExpirationDate  time.Time                 `gorm:"not null" sort:"expirationDate"`
}

// Hooks for SupplierWarehouseDeliveryMaterial
func (s *SupplierWarehouseDeliveryMaterial) BeforeCreate(tx *gorm.DB) (err error) {
	s.ExpirationDate = toUTC(s.ExpirationDate)
	return
}

func (s *SupplierWarehouseDeliveryMaterial) AfterFind(tx *gorm.DB) (err error) {
	s.ExpirationDate = toUTC(s.ExpirationDate)
	return
}

type AggregatedWarehouseStock struct {
	WarehouseID            uint          `json:"warehouseId"`
	StockMaterialID        uint          `json:"stockMaterialId"`
	StockMaterial          StockMaterial `gorm:"foreignKey:StockMaterialID;references:ID;preload:true" json:"stockMaterial"`
	TotalQuantity          float64       `json:"totalQuantity"`
	EarliestExpirationDate *time.Time    `json:"earliestExpirationDate"`
}
