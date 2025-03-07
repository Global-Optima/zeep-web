package data

type OrderStatus string

const (
	OrderStatusWaitingForPayment OrderStatus = "WAITING_FOR_PAYMENT"
	OrderStatusPending           OrderStatus = "PENDING"
	OrderStatusPreparing         OrderStatus = "PREPARING"
	OrderStatusCompleted         OrderStatus = "COMPLETED"
	OrderStatusInDelivery        OrderStatus = "IN_DELIVERY"
	OrderStatusDelivered         OrderStatus = "DELIVERED"
	OrderStatusCancelled         OrderStatus = "CANCELLED"
)

type SubOrderStatus string

const (
	SubOrderStatusPending   SubOrderStatus = "PENDING"
	SubOrderStatusPreparing SubOrderStatus = "PREPARING"
	SubOrderStatusCompleted SubOrderStatus = "COMPLETED"
)

type TransactionType string

const (
	TransactionTypePayment TransactionType = "PAYMENT"
	TransactionTypeRefund  TransactionType = "REFUND"
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
	Transactions      []Transaction   `gorm:"foreignKey:OrderID;constraint:OnDelete:SET NULL"`
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

type Transaction struct {
	BaseEntity
	Type          TransactionType `gorm:"type:varchar(255);not null" sort:"type"`
	OrderID       uint            `gorm:"index;not null"`
	Order         Order           `gorm:"foreignKey:OrderID;constraint:OnUpdate:SET NULL"`
	Bin           string          `gorm:"type:varchar(20);not null"`
	TransactionID string          `gorm:"type:varchar(20);unique;not null"`
	ProcessID     *string         `gorm:"type:varchar(20);unique;nullable"`
	PaymentMethod string          `gorm:"type:varchar(50);not null"`
	Amount        float64         `gorm:"type:decimal(10,2);not null" sort:"amount"`
	Currency      string          `gorm:"type:char(3);not null"`
	QRNumber      *string         `gorm:"type:varchar(16)"`
	CardMask      *string         `gorm:"type:varchar(16)"`
	ICC           *string         `gorm:"type:varchar(255)"`
}
