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
	CustomerName      string              `json:"customerName" binding:"required,regexp=^[a-zA-ZәӘіІңҢғҒүҮұҰқҚөӨһҺ\\s]+$"`
	EmployeeID        *uint               `json:"employeeId,omitempty"`
	DeliveryAddressID *uint               `json:"deliveryAddressId"`
	Suborders         []CreateSubOrderDTO `json:"subOrders"`

	StoreID uint
}

type OrderStatusesCountDTO struct {
	ALL         int64 `json:"ALL"`
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
	EmployeeID        *uint            `json:"employeeId,omitempty"`
	StoreID           uint             `json:"storeId"`
	DeliveryAddressID *uint            `json:"deliveryAddressId,omitempty"`
	Status            data.OrderStatus `json:"status"`
	CreatedAt         time.Time        `json:"createdAt"`
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
}

type OrderStoreProductSizeDTO struct {
	ID          uint               `json:"id"`
	SizeName    string             `json:"sizeName"`
	ProductName string             `json:"productName"`
	Size        int                `json:"size"`
	Unit        unitTypes.UnitsDTO `json:"unit"`
}

type OrderStoreAdditiveDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
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
