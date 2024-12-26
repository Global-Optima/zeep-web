package data

import "time"

type StockRequestStatus string

var (
	StockRequestCreated    StockRequestStatus = "CREATED"
	StockRequestInDelivery StockRequestStatus = "IN_DELIVERY"
	StockRequestCompleted  StockRequestStatus = "COMPLETED"
)

type Warehouse struct {
	BaseEntity
	FacilityAddressID uint            `gorm:"not null;index"`
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
	LowStockThreshold float64        `gorm:"type:decimal(10,2);not null;check:quantity > 0" sort:"lowStockThreshold"`
	Quantity          float64        `gorm:"type:decimal(10,2);not null;check:quantity >= 0" sort:"quantity"`
}

type StockRequest struct {
	BaseEntity
	WarehouseID uint                     `gorm:"not null;index"`
	Warehouse   Warehouse                `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE" sort:"warehouses"`
	Status      StockRequestStatus       `gorm:"size:50;not null" sort:"status"`
	RequestDate *time.Time               `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	Ingredients []StockRequestIngredient `gorm:"foreignKey:StockRequestID;constraint:OnDelete:CASCADE"`
}

type StockRequestIngredient struct {
	BaseEntity
	StockRequestID uint         `gorm:"not null;index"`
	StockRequest   StockRequest `gorm:"foreignKey:StockRequestID;constraint:OnDelete:CASCADE"`
	IngredientID   uint         `gorm:"not null;index"`
	Ingredient     Ingredient   `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	Quantity       float64      `gorm:"type:decimal(10,2);not null;check:quantity > 0" sort:"quantity"`
	DeliveredDate  time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" sort:"deliveryDate"` // Delivery start date
	ExpirationDate time.Time    `gorm:"not null" sort:"expirationDate"`                         // Calculated from DeliveredDate + ExpirationPeriod
}

type IngredientsMapping struct {
	BaseEntity
	IngredientID    uint          `gorm:"not null;index"`
	Ingredient      Ingredient    `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	StockMaterialID uint          `gorm:"not null;index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
}

type Supplier struct {
	BaseEntity
	Name         string `gorm:"size:255;not null" sort:"name"`
	ContactEmail string `gorm:"size:255"`
	ContactPhone string `gorm:"size:16"`
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
	Name                   string   `gorm:"size:255;not null" sort:"name"`
	Description            string   `gorm:"type:text"`
	SafetyStock            float64  `gorm:"type:decimal(10,2);not null" sort:"safetyStock"`
	ExpirationFlag         bool     `gorm:"not null" sort:"expirationFlag"`
	UnitID                 uint     `gorm:"not null"`
	Unit                   Unit     `gorm:"foreignKey:UnitID;constraint:OnDelete:SET NULL"`
	Category               string   `gorm:"size:255" sort:"category"`
	Barcode                string   `gorm:"unique;size:255"`
	ExpirationPeriodInDays int      `gorm:"not null;default:1095" sort:"expirationPeriodInDays"` // 3 years in days
	IsActive               bool     `gorm:"not null;default:true" sort:"isActive"`
	Package                *Package `gorm:"foreignKey:StockMaterialID"`
}

type Unit struct {
	BaseEntity
	Name             string  `gorm:"size:50;not null;unique"`
	ConversionFactor float64 `gorm:"type:decimal(10,4);not null"` // To base unit
}

type Package struct {
	BaseEntity
	StockMaterialID uint          `gorm:"index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	PackageSize     float64       `gorm:"type:decimal(10,2);not null"`
	PackageUnitID   uint          `gorm:"not null"`
	PackageUnit     Unit          `gorm:"foreignKey:PackageUnitID;constraint:OnDelete:SET NULL"`
}

type SupplierMaterial struct {
	BaseEntity
	StockMaterialID uint          `gorm:"not null;index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	SupplierID      uint          `gorm:"not null;index"`
	Supplier        Supplier      `gorm:"foreignKey:SupplierID;constraint:OnDelete:CASCADE"`
}

type Delivery struct {
	BaseEntity
	StockMaterialID uint          `gorm:"not null;index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	SupplierID      uint          `gorm:"not null"`
	WarehouseID     uint          `gorm:"not null"`
	Barcode         string        `gorm:"size:255;not null"`
	Quantity        float64       `gorm:"type:decimal(10,2);not null;check:quantity > 0" sort:"quantity"`
	DeliveryDate    time.Time     `gorm:"not null" sort:"deliveryDate"`
	ExpirationDate  time.Time     `gorm:"not null" sort:"expirationDate"`
}
