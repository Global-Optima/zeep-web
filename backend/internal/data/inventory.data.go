package data

import "time"

type StockRequestStatus string

var (
	StockRequestCreated    StockRequestStatus = "CREATED"
	StockRequestProcessed  StockRequestStatus = "PROCESSED"
	StockRequestInDelivery StockRequestStatus = "IN_DELIVERY"
	StockRequestCompleted  StockRequestStatus = "COMPLETED"
	StockRequestRejected   StockRequestStatus = "REJECTED"
)

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
	StoreID     uint      `gorm:"not null;index"`
	Store       Store     `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	WarehouseID uint      `gorm:"not null;index"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE"`
}

type StoreWarehouseStock struct {
	BaseEntity
	StoreWarehouseID  uint           `gorm:"not null;index"`
	StoreWarehouse    StoreWarehouse `gorm:"foreignKey:StoreWarehouseID;constraint:OnDelete:CASCADE"`
	IngredientID      uint           `gorm:"not null;index"`
	Ingredient        Ingredient     `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	LowStockThreshold float64        `gorm:"type:decimal(10,2);not null;check:quantity > 0"`
	Quantity          float64        `gorm:"type:decimal(10,2);not null;check:quantity >= 0"`
}

type StockRequest struct {
	BaseEntity
	StoreID     uint                     `gorm:"not null;index"` // Links to Store
	Store       Store                    `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	WarehouseID uint                     `gorm:"not null;index"` // Central warehouse fulfilling the request
	Warehouse   Warehouse                `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE"`
	Status      StockRequestStatus       `gorm:"size:50;not null"`
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
	DeliveredDate  time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP"` // Delivery start date
	ExpirationDate time.Time    `gorm:"not null"`                           // Calculated from DeliveredDate + ExpirationPeriodInDays
}

type IngredientStockMaterialMapping struct {
	BaseEntity
	IngredientID    uint          `gorm:"not null;index"`
	Ingredient      Ingredient    `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	StockMaterialID uint          `gorm:"not null;index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
}

type Supplier struct {
	BaseEntity
	Name         string `gorm:"size:255;not null"`
	ContactEmail string `gorm:"size:255"`
	ContactPhone string `gorm:"size:20"`
	Address      string `gorm:"size:255"`
}

type WarehouseStock struct {
	BaseEntity
	WarehouseID     uint          `gorm:"not null;index"`
	Warehouse       Warehouse     `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE"`
	StockMaterialID uint          `gorm:"index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	Quantity        float64       `gorm:"type:decimal(10,2);not null;check:quantity >= 0"`
}

type StockMaterial struct {
	BaseEntity
	Name                   string                `gorm:"size:255;not null"`
	Description            string                `gorm:"type:text"`
	SafetyStock            float64               `gorm:"type:decimal(10,2);not null"`
	ExpirationFlag         bool                  `gorm:"not null"`
	UnitID                 uint                  `gorm:"not null"`
	Unit                   Unit                  `gorm:"foreignKey:UnitID;constraint:OnDelete:SET NULL"`
	Category               string                `gorm:"size:255"`
	Barcode                string                `gorm:"unique;size:255"`
	ExpirationPeriodInDays int                   `gorm:"not null;default:1095"` // 3 years in days
	IsActive               bool                  `gorm:"not null;default:true"`
	Package                *StockMaterialPackage `gorm:"foreignKey:StockMaterialID"`
}

type Unit struct {
	BaseEntity
	Name             string  `gorm:"size:50;not null;unique"`
	ConversionFactor float64 `gorm:"type:decimal(10,4);not null"` // To base unit
}

type StockMaterialPackage struct {
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

type SupplierWarehouseDelivery struct {
	BaseEntity
	StockMaterialID uint          `gorm:"not null;index"`
	StockMaterial   StockMaterial `gorm:"foreignKey:StockMaterialID;constraint:OnDelete:CASCADE"`
	SupplierID      uint          `gorm:"not null"`
	WarehouseID     uint          `gorm:"not null"`
	Barcode         string        `gorm:"size:255;not null"`
	Quantity        float64       `gorm:"type:decimal(10,2);not null;check:quantity > 0"`
	DeliveryDate    time.Time     `gorm:"not null"`
	ExpirationDate  time.Time     `gorm:"not null"`
}
