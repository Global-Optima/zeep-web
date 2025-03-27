package data

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Franchisee struct {
	BaseEntity
	Name                string               `gorm:"size:255;not null" sort:"name"`
	Description         string               `gorm:"size:1024"`
	FranchiseeEmployees []FranchiseeEmployee `gorm:"foreignKey:FranchiseeID"`
	Stores              []Store              `gorm:"foreignKey:FranchiseeID"`
}

type Store struct {
	BaseEntity
	Name                string          `gorm:"size:255;not null" sort:"name"`
	FacilityAddressID   uint            `gorm:"index;not null"`
	FacilityAddress     FacilityAddress `gorm:"foreignKey:FacilityAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	FranchiseeID        *uint           `gorm:"index"`
	Franchisee          *Franchisee     `gorm:"foreignKey:FranchiseeID" sort:"franchisees"`
	WarehouseID         uint            `gorm:"not null;index"` // New Warehouse Reference
	Warehouse           Warehouse       `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE"`
	IsActive            bool            `gorm:"default:true" sort:"isActive"`
	ContactPhone        string          `gorm:"size:16"`
	ContactEmail        string          `gorm:"size:255"`
	StoreHours          string          `gorm:"size:255"`
	Additives           []StoreAdditive `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Products            []StoreProduct  `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Stocks              []StoreStock    `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"` // Linked Store Stocks
	LastInventorySyncAt time.Time       `gorm:"autoCreateTime" sort:"lastInventorySyncAt"`
}

type StoreStock struct {
	BaseEntity
	StoreID           uint       `gorm:"not null;index"`
	Store             Store      `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	IngredientID      uint       `gorm:"not null;index"`
	Ingredient        Ingredient `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE" sort:"ingredients"`
	LowStockThreshold float64    `gorm:"type:decimal(10,2);not null;check:low_stock_threshold > 0" sort:"lowStockThreshold"`
	Quantity          float64    `gorm:"type:decimal(10,2);not null;check:quantity >= 0" sort:"quantity"`
}

func (storeStock *StoreStock) AfterUpdate(tx *gorm.DB) error {

	if storeStock == nil || storeStock.ID == 0 || storeStock.StoreID == 0 || storeStock.IngredientID == 0 {
		return fmt.Errorf("not enough info to fire AfterUpdate hook, preload the entity to proceed or use save instead of updates")
	}

	return UpdateOutOfStockInventory(tx, storeStock.StoreID, storeStock.IngredientID)
}

func (storeStock *StoreStock) AfterCreate(tx *gorm.DB) error {
	return UpdateOutOfStockInventory(tx, storeStock.StoreID, storeStock.IngredientID)
}

type StoreAdditive struct {
	BaseEntity
	AdditiveID   uint     `gorm:"index;not null"`
	StoreID      uint     `gorm:"index;not null"`
	StorePrice   *float64 `gorm:"type:decimal(10,2);check:price >= 0"`
	IsOutOfStock bool     `gorm:"not null;default:false" sort:"isOutOfStock"`
	Store        Store    `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Additive     Additive `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
}

type StoreProductSize struct {
	BaseEntity
	ProductSizeID  uint         `gorm:"index;not null"`
	StoreProductID uint         `gorm:"index;not null"`
	StorePrice     *float64     `gorm:"type:decimal(10,2);not null;check:price >= 0"`
	StoreProduct   StoreProduct `gorm:"foreignKey:StoreProductID;constraint:OnDelete:CASCADE"`
	ProductSize    ProductSize  `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
}

type StoreProduct struct {
	BaseEntity
	ProductID         uint               `gorm:"index;not null"`
	StoreID           uint               `gorm:"index;not null"`
	IsAvailable       bool               `gorm:"default:true" sort:"isAvailable"`
	IsOutOfStock      bool               `gorm:"not null;default:false" sort:"isOutOfStock"`
	Store             Store              `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Product           Product            `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" sort:"products"`
	StoreProductSizes []StoreProductSize `gorm:"foreignKey:StoreProductID;constraint:OnDelete:CASCADE"`
}

type FacilityAddress struct {
	BaseEntity
	Address   string   `gorm:"size:255;not null"`
	Longitude *float64 `gorm:"type:decimal(9,6)"`
	Latitude  *float64 `gorm:"type:decimal(9,6)"`
	Stores    []Store  `gorm:"foreignKey:FacilityAddressID"`
}
