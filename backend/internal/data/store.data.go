package data

import (
	"time"
)

type ProvisionStatus string

const (
	PROVISION_STATUS_PREPARING ProvisionStatus = "PREPARING"
	PROVISION_STATUS_COMPLETED ProvisionStatus = "COMPLETED"
	PROVISION_STATUS_EXPIRED   ProvisionStatus = "EXPIRED"
)

type Product struct {
	BaseEntity
	Name         string           `gorm:"size:100;not null" sort:"name"`
	Description  string           `gorm:"type:text"`
	ImageKey     *StorageImageKey `gorm:"size:2048"`
	VideoKey     *StorageVideoKey `gorm:"size:2048"`
	CategoryID   uint             `gorm:"index;not null"`
	Category     ProductCategory  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" sort:"category"`
	RecipeSteps  []RecipeStep     `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	ProductSizes []ProductSize    `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

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

type StoreProvision struct {
	BaseEntity
	ProvisionID uint            `gorm:"not null;index"`
	Provision   Provision       `gorm:"foreignKey:ProvisionID;constraint:OnDelete:CASCADE"`
	Volume      float64         `gorm:"type:decimal(10,2);not null;check:volume > 0" sort:"volume"`
	Status      ProvisionStatus `gorm:"not null" sort:"status"`
	CompletedAt time.Time       `gorm:"" sort:"completedAt"`
	ExpiresAt   time.Time       `gorm:"not null" sort:"expiresAt"`
}

type StoreProvisionIngredient struct {
	BaseEntity
	IngredientID     uint           `gorm:"not null;index"`
	Ingredient       Ingredient     `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	StoreProvisionID uint           `gorm:"not null"`
	StoreProvision   StoreProvision `gorm:"foreignKey:StoreProvisionID;constraint:OnDelete:CASCADE"`
	Quantity         float64        `gorm:"not null;check:quantity > 0"`
}
