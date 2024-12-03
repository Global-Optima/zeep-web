package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type CreateOrderDTO struct {
	CustomerID        *uint                `json:"customerId,omitempty"`
	CustomerName      string               `json:"customerName"`
	EmployeeID        *uint                `json:"employeeId,omitempty"`
	StoreID           uint                 `json:"storeId"`
	DeliveryAddressID *uint                `json:"deliveryAddressId"`
	OrderItems        []CreateOrderItemDTO `json:"orderItems"`
	OrderType         string               `json:"type"` // delivery or in cafe
}

type CreateOrderItemDTO struct {
	ProductSizeID uint   `json:"productSizeId"`
	Quantity      int    `json:"quantity"`
	AdditivesIDs  []uint `json:"additivesIds"`
}

type OrderDTO struct {
	ID                uint              `json:"id"`
	CustomerID        *uint             `json:"customerId"`
	CustomerName      *string           `json:"customerName"`
	EmployeeID        *uint             `json:"employeeId,omitempty"`
	StoreID           uint              `json:"storeId"`
	DeliveryAddressID *uint             `json:"deliveryAddressId,omitempty"`
	OrderStatus       data.OrderStatus  `json:"orderStatus"`
	CreatedAt         time.Time         `json:"orderDate"`
	Total             float64           `json:"total"`
	SubOrdersQuantity int               `json:"sub_orders_quantity"`
	OrderProducts     []OrderProductDTO `json:"orderProducts,omitempty"`
	Timestamp         time.Time         `json:"timestamp"`
}

type OrderProductDTO struct {
	ID            uint                      `json:"id"`
	OrderID       uint                      `json:"orderId"`
	ProductSizeID uint                      `json:"productSizeId"`
	Quantity      int                       `json:"quantity"`
	Price         float64                   `json:"price"`
	Additives     []OrderProductAdditiveDTO `json:"additives,omitempty"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
}

type OrderProductAdditiveDTO struct {
	ID             uint      `json:"id"`
	OrderProductID uint      `json:"orderProductId"`
	AdditiveID     uint      `json:"additiveId"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
