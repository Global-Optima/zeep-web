package data

import "time"

type Order struct {
	BaseEntity
	CustomerID        uint           `gorm:"index;not null"`
	EmployeeID        *uint          `gorm:"index"`
	StoreID           *uint          `gorm:"index"`
	DeliveryAddressID *uint          `gorm:"index"`
	OrderStatus       string         `gorm:"size:50;not null"` // e.g., 'pending', 'completed'
	OrderDate         time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	Total             float64        `gorm:"type:decimal(10,2);not null;check:total >= 0"`
	OrderProducts     []OrderProduct `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderProduct struct {
	BaseEntity
	OrderID       uint                   `gorm:"index;not null"`
	ProductSizeID uint                   `gorm:"index;not null"` // Links to ProductSize
	ProductSize   ProductSize            `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	Quantity      int                    `gorm:"not null;check:quantity > 0"`
	Price         float64                `gorm:"type:decimal(10,2);not null"`
	Additives     []OrderProductAdditive `gorm:"foreignKey:OrderProductID;constraint:OnDelete:CASCADE"`
}

type OrderProductAdditive struct {
	BaseEntity
	OrderProductID uint     `gorm:"index;not null"`
	AdditiveID     uint     `gorm:"index;not null"` // Links to Additive
	Additive       Additive `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
	Price          float64  `gorm:"type:decimal(10,2);not null"`
}
