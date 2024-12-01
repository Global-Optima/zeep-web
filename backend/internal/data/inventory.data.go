package data

import "time"

type CityWarehouse struct {
	BaseEntity
	FacilityAddressID uint            `gorm:"not null;index"`
	FacilityAddress   FacilityAddress `gorm:"foreignKey:FacilityAddressID;constraint:OnDelete:CASCADE"`
	Name              string          `gorm:"size:255;not null"`
}

type StoreWarehouse struct {
	BaseEntity
	StoreID         uint          `gorm:"not null;index"`
	Store           Store         `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	CityWarehouseID uint          `gorm:"not null;index"`
	CityWarehouse   CityWarehouse `gorm:"foreignKey:CityWarehouseID;constraint:OnDelete:CASCADE"`
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
	CityWarehouseID uint                     `gorm:"not null;index"`
	CityWarehouse   CityWarehouse            `gorm:"foreignKey:CityWarehouseID;constraint:OnDelete:CASCADE"`
	Status          string                   `gorm:"size:50;not null"`
	RequestDate     *time.Time               `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	Ingredients     []StockRequestIngredient `gorm:"foreignKey:StockRequestID;constraint:OnDelete:CASCADE"`
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
