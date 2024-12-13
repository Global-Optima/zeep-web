package data

type OrderStatus string

const (
	OrderStatusPreparing  OrderStatus = "PREPARING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusInDelivery OrderStatus = "IN_DELIVERY"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type SubOrderStatus string

const (
	SubOrderStatusPreparing SubOrderStatus = "PREPARING"
	SubOrderStatusCompleted SubOrderStatus = "COMPLETED"
)

type Order struct {
	BaseEntity
	CustomerID        *uint       `gorm:"index"`
	CustomerName      string      `gorm:"size:255"`
	EmployeeID        *uint       `gorm:"index"`
	StoreID           uint        `gorm:"index"`
	DeliveryAddressID *uint       `gorm:"index"`
	Status            OrderStatus `gorm:"size:50;not null"`
	Total             float64     `gorm:"type:decimal(10,2);not null;check:total >= 0"`
	Suborders         []Suborder  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

// Suborder Model
type Suborder struct {
	BaseEntity
	OrderID       uint               `gorm:"index;not null"`
	ProductSizeID uint               `gorm:"index;not null"`
	ProductSize   ProductSize        `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	Price         float64            `gorm:"type:decimal(10,2);not null;check:price >= 0"`
	Status        SubOrderStatus     `gorm:"size:50;not null"`
	Additives     []SuborderAdditive `gorm:"foreignKey:SuborderID;constraint:OnDelete:CASCADE"`
}

// SuborderAdditive Model
type SuborderAdditive struct {
	BaseEntity
	SuborderID uint     `gorm:"index;not null"`
	AdditiveID uint     `gorm:"index;not null"`
	Additive   Additive `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
	Price      float64  `gorm:"type:decimal(10,2);not null;check:price >= 0"`
}
