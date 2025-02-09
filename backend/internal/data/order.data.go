package data

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusPreparing  OrderStatus = "PREPARING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusInDelivery OrderStatus = "IN_DELIVERY"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type SubOrderStatus string

const (
	SubOrderStatusPending   SubOrderStatus = "PENDING"
	SubOrderStatusPreparing SubOrderStatus = "PREPARING"
	SubOrderStatusCompleted SubOrderStatus = "COMPLETED"
)

type Order struct {
	BaseEntity
	CustomerID        *uint           `gorm:"index"`
	CustomerName      string          `gorm:"size:255" sort:"customerName"`
	StoreEmployeeID   *uint           `gorm:"index"`
	StoreEmployee     *StoreEmployee  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StoreID           uint            `gorm:"index"`
	Store             Store           `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	DeliveryAddressID *uint           `gorm:"index"`
	DeliveryAddress   CustomerAddress `gorm:"foreignKey:DeliveryAddressID;constraint:OnDelete:CASCADE"`
	Status            OrderStatus     `gorm:"size:50;not null" sort:"orderStatus"`
	Total             float64         `gorm:"type:decimal(10,2);not null;check:total >= 0" sort:"total"`
	Suborders         []Suborder      `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	DisplayNumber     int             `gorm:"not null;index"`
}

// Suborder Model
type Suborder struct {
	BaseEntity
	OrderID            uint               `gorm:"index;not null"`
	StoreProductSizeID uint               `gorm:"index;not null"`
	StoreProductSize   StoreProductSize   `gorm:"foreignKey:StoreProductSizeID;constraint:OnDelete:CASCADE"`
	Price              float64            `gorm:"type:decimal(10,2);not null;check:price >= 0"`
	Status             SubOrderStatus     `gorm:"size:50;not null"`
	SuborderAdditives  []SuborderAdditive `gorm:"foreignKey:SuborderID;constraint:OnDelete:CASCADE"`
}

// SuborderAdditive Model
type SuborderAdditive struct {
	BaseEntity
	SuborderID      uint          `gorm:"index;not null"`
	StoreAdditiveID uint          `gorm:"index;not null"`
	StoreAdditive   StoreAdditive `gorm:"foreignKey:StoreAdditiveID;constraint:OnDelete:CASCADE"`
	Price           float64       `gorm:"type:decimal(10,2);not null;check:price >= 0"`
}
