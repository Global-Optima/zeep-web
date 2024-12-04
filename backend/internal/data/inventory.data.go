package data

import "time"

type Warehouse struct {
	BaseEntity
	FacilityAddressID uint            `gorm:"not null;index"`
	FacilityAddress   FacilityAddress `gorm:"foreignKey:FacilityAddressID;constraint:OnDelete:CASCADE"`
	Name              string          `gorm:"size:255;not null"`
}

type StoreWarehouse struct {
	BaseEntity
	StoreID     uint      `gorm:"not null;index"`
	Store       Store     `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	WarehouseID uint      `gorm:"not null;index"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE"`
}

type StoreWarehouseStock struct {
	BaseEntity
	StoreWarehouseID uint           `gorm:"not null;index"`
	StoreWarehouse   StoreWarehouse `gorm:"foreignKey:StoreWarehouseID;constraint:OnDelete:CASCADE"`
	IngredientID     uint           `gorm:"not null;index"`
	Ingredient       Ingredient     `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	Quantity         float64        `gorm:"type:decimal(10,2);not null;check:quantity >= 0"`
}

type StockRequest struct {
	BaseEntity
	WarehouseID uint                     `gorm:"not null;index"`
	Warehouse   Warehouse                `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE"`
	Status      string                   `gorm:"size:50;not null"`
	RequestDate *time.Time               `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	Ingredients []StockRequestIngredient `gorm:"foreignKey:StockRequestID;constraint:OnDelete:CASCADE"`
}

type StockRequestIngredient struct {
	BaseEntity
	StockRequestID uint         `gorm:"not null;index"`
	StockRequest   StockRequest `gorm:"foreignKey:StockRequestID;constraint:OnDelete:CASCADE"`
	IngredientID   uint         `gorm:"not null;index"`
	Ingredient     Ingredient   `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	Quantity       float64      `gorm:"type:decimal(10,2);not null;check:quantity > 0"`
}

type Supplier struct {
	BaseEntity
	Name         string `gorm:"size:255;not null"`
	ContactEmail string `gorm:"size:255"`
	ContactPhone string `gorm:"size:20"`
	Address      string `gorm:"size:255"`
}

type SKU struct {
	BaseEntity
	ID               string  `gorm:"unique;not null"`
	Name             string  `gorm:"size:255;not null"`
	Description      string  `gorm:"type:text"`
	SafetyStock      float64 `gorm:"type:decimal(10,2);not null"`
	ExpirationFlag   bool    `gorm:"not null"`
	Quantity         float64 `gorm:"type:decimal(10,2);not null;check:quantity >= 0"`
	UnitID           uint    `gorm:"not null"`
	Unit             Unit    `gorm:"foreignKey:UnitID;constraint:OnDelete:SET NULL"`
	Category         string  `gorm:"size:255"`
	Barcode          string  `gorm:"unique;size:255"`
	ExpirationPeriod int     `gorm:"not null;default:1095"` // 3 years
	IsActive         bool    `gorm:"not null;default:true"`
}

type Unit struct {
	BaseEntity
	Name             string  `gorm:"size:50;not null;unique"`
	ConversionFactor float64 `gorm:"type:decimal(10,4);not null"` // To base unit
}

type Package struct {
	BaseEntity
	SKU_ID        uint    `gorm:"not null;index"`
	SKU           SKU     `gorm:"foreignKey:SKU_ID;constraint:OnDelete:CASCADE"`
	PackageSize   float64 `gorm:"type:decimal(10,2);not null"`
	PackageUnitID uint    `gorm:"not null"`
	PackageUnit   Unit    `gorm:"foreignKey:PackageUnitID;constraint:OnDelete:SET NULL"`
}
