package types

import (
	"time"

	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type OrdersFilterQuery struct {
	Search  *string           `form:"search"`
	Status  *data.OrderStatus `form:"status"`
	StoreID *uint             `form:"storeId"`
	utils.BaseFilter
}

type CreateOrderDTO struct {
	CustomerID        *uint               `json:"customerId,omitempty"`
	CustomerName      string              `json:"customerName" binding:"required"`
	StoreEmployeeID   *uint               `json:"storeEmployeeId,omitempty"`
	DeliveryAddressID *uint               `json:"deliveryAddressId"`
	Suborders         []CreateSubOrderDTO `json:"subOrders"`

	StoreID uint
}

type ValidateCustomerNameDTO struct {
	CustomerName string `json:"customerName" binding:"required"`
}

type OrderStatusesCountDTO struct {
	ALL         int64 `json:"ALL"`
	PENDING     int64 `json:"PENDING"`
	PREPARING   int64 `json:"PREPARING"`
	COMPLETED   int64 `json:"COMPLETED"`
	IN_DELIVERY int64 `json:"IN_DELIVERY"`
	DELIVERED   int64 `json:"DELIVERED"`
	CANCELLED   int64 `json:"CANCELLED"`
}

type CreateSubOrderDTO struct {
	StoreProductSizeID uint   `json:"storeProductSizeId"`
	Quantity           int    `json:"quantity"`
	StoreAdditivesIDs  []uint `json:"storeAdditivesIds"`
}

type OrderDTO struct {
	ID                uint             `json:"id"`
	CustomerID        *uint            `json:"customerId,omitempty"`
	CustomerName      *string          `json:"customerName"`
	StoreEmployeeID   *uint            `json:"storeEmployeeId,omitempty"`
	StoreID           uint             `json:"storeId"`
	DeliveryAddressID *uint            `json:"deliveryAddressId,omitempty"`
	Status            data.OrderStatus `json:"status"`
	CreatedAt         time.Time        `json:"createdAt"`
	CompletedAt       *time.Time       `json:"completedAt,omitempty"`
	Total             float64          `json:"total"`
	DisplayNumber     int              `json:"displayNumber"`
	SubordersQuantity int              `json:"subOrdersQuantity"`
	Suborders         []SuborderDTO    `json:"subOrders"`
}

type SuborderDTO struct {
	ID          uint                       `json:"id"`
	OrderID     uint                       `json:"orderId"`
	ProductSize OrderStoreProductSizeDTO   `json:"productSize"`
	Price       float64                    `json:"price"`
	Status      data.SubOrderStatus        `json:"status"`
	Additives   []SuborderStoreAdditiveDTO `json:"additives"`
	CreatedAt   time.Time                  `json:"createdAt"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
	CompletedAt *time.Time                 `json:"completedAt,omitempty"`
}

type OrderStoreProductSizeDTO struct {
	ID          uint               `json:"id"`
	SizeName    string             `json:"sizeName"`
	ProductName string             `json:"productName"`
	Size        float64            `json:"size"`
	Unit        unitTypes.UnitsDTO `json:"unit"`
	MachineId   string             `json:"machineId"`
}

type OrderStoreAdditiveDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Size        float64 `json:"size"`
	MachineId   string  `json:"machineId"`
}

type SuborderStoreAdditiveDTO struct {
	ID         uint                  `json:"id"`
	SuborderID uint                  `json:"subOrderId"`
	Additive   OrderStoreAdditiveDTO `json:"additive"`
	Price      float64               `json:"price"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

type OrderDetailsDTO struct {
	ID              uint                     `json:"id"`
	CustomerName    *string                  `json:"customerName,omitempty"` // Optional
	Status          string                   `json:"status"`
	Total           float64                  `json:"total"`
	Suborders       []SuborderDetailsDTO     `json:"suborders"`
	DeliveryAddress *OrderDeliveryAddressDTO `json:"deliveryAddress,omitempty"` // Optional
}

type SuborderDetailsDTO struct {
	ID               uint                       `json:"id"`
	Price            float64                    `json:"price"`
	Status           string                     `json:"status"`
	StoreProductSize OrderProductSizeDetailsDTO `json:"storeProductSize"`
	StoreAdditives   []OrderAdditiveDetailsDTO  `json:"storeAdditives"`
}

type OrderProductSizeDetailsDTO struct {
	ID         uint                   `json:"id"`
	Name       string                 `json:"name"`
	Unit       unitTypes.UnitsDTO     `json:"unit"`
	StorePrice float64                `json:"storePrice"`
	Product    OrderProductDetailsDTO `json:"product"`
}

type OrderProductDetailsDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type OrderAdditiveDetailsDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"basePrice"`
}

type OrderDeliveryAddressDTO struct {
	ID        uint   `json:"id"`
	Address   string `json:"address"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type OrdersExportFilterQuery struct {
	StartDate *time.Time `form:"startDate" binding:"omitempty"`
	EndDate   *time.Time `form:"endDate" binding:"omitempty"`
	StoreID   *uint      `form:"storeId" binding:"omitempty"`
	Language  string     `form:"language" binding:"omitempty,oneof=kk ru en"` // Optional language filter
}

type OrderExportDTO struct {
	ID              uint                     `json:"id"`
	CustomerName    string                   `json:"customerName"`
	Status          string                   `json:"status"`
	Total           float64                  `json:"total"`
	CreatedAt       time.Time                `json:"createdAt"`
	StoreName       string                   `json:"storeName"`
	Suborders       []SuborderDTO            `json:"suborders"`
	DeliveryAddress *OrderDeliveryAddressDTO `json:"deliveryAddress,omitempty"`
}

type OrdersTimeZoneFilter struct {
	StoreID          *uint   `form:"storeId" binding:"omitempty"`
	TimeZoneLocation *string `form:"timezone" binding:"omitempty"`
	TimeZoneOffset   *uint   `form:"timezoneOffset" binding:"omitempty"`
}

type TransactionDTO struct {
	Bin           string  `json:"bin" binding:"required,len=12"`
	TransactionID string  `json:"transactionId" binding:"required,max=20"`
	ProcessID     *string `json:"processId" binding:"omitempty,max=20"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required,len=3"`
	QRNumber      *string `json:"qrNumber" binding:"omitempty,min=7,max=16"`
	CardMask      *string `json:"cardMask" binding:"omitempty,len=16"`
	ICC           *string `json:"icc" binding:"omitempty"`
}

type WaitingOrderPayload struct {
	OrderID uint `json:"orderId"`
}
