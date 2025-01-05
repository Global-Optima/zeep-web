package types

import (
	"time"

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
	CustomerName      string              `json:"customerName"`
	EmployeeID        *uint               `json:"employeeId,omitempty"`
	StoreID           uint                `json:"storeId"`
	DeliveryAddressID *uint               `json:"deliveryAddressId"`
	Suborders         []CreateSubOrderDTO `json:"subOrders"`
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
	ProductSizeID uint   `json:"productSizeId"`
	Quantity      int    `json:"quantity"`
	AdditivesIDs  []uint `json:"additivesIds"`
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
	SubordersQuantity int              `json:"subOrdersQuantity"`
	Suborders         []SuborderDTO    `json:"subOrders"`
}

type SuborderDTO struct {
	ID          uint                  `json:"id"`
	OrderID     uint                  `json:"orderId"`
	ProductSize ProductSizeDTO        `json:"productSize"`
	Price       float64               `json:"price"`
	Status      data.SubOrderStatus   `json:"status"`
	Additives   []SuborderAdditiveDTO `json:"additives"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

type ProductSizeDTO struct {
	ID          uint   `json:"id"`
	SizeName    string `json:"sizeName"`
	ProductName string `json:"productName"`
	Size        int    `json:"size"`
}

type AdditiveDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        string `json:"size"`
}

type SuborderAdditiveDTO struct {
	ID         uint        `json:"id"`
	SuborderID uint        `json:"subOrderId"`
	Additive   AdditiveDTO `json:"additive"`
	Price      float64     `json:"price"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type OrderDetailsDTO struct {
	ID              uint                 `json:"id"`
	CustomerName    *string              `json:"customerName,omitempty"` // Optional
	Status          string               `json:"status"`
	Total           float64              `json:"total"`
	Suborders       []SuborderDetailsDTO `json:"suborders"`
	DeliveryAddress *DeliveryAddressDTO  `json:"deliveryAddress,omitempty"` // Optional
}

type SuborderDetailsDTO struct {
	ID          uint                  `json:"id"`
	Price       float64               `json:"price"`
	Status      string                `json:"status"`
	ProductSize ProductSizeDetailsDTO `json:"productSize"`
	Additives   []AdditiveDetailsDTO  `json:"additives"`
}

type ProductSizeDetailsDTO struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Measure   string            `json:"measure"`
	BasePrice float64           `json:"basePrice"`
	Product   ProductDetailsDTO `json:"product"`
}

type ProductDetailsDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type AdditiveDetailsDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"basePrice"`
}

type DeliveryAddressDTO struct {
	ID        uint   `json:"id"`
	Address   string `json:"address"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
