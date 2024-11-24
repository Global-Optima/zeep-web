package types

import (
	"time"
)

type CreateOrderDTO struct {
	CustomerID        uint                `json:"customer_id"`
	EmployeeID        *uint               `json:"employee_id,omitempty"`
	StoreID           *uint               `json:"store_id"`
	DeliveryAddressID *uint               `json:"delivery_address_id"`
	Products          []CreateSubOrderDTO `json:"products"`
}

type CreateSubOrderDTO struct {
	ProductSizeID uint   `json:"product_size_id"`
	Quantity      int    `json:"quantity"`
	Additives     []uint `json:"additives"` //list of ids
}

type OrderDTO struct {
	ID                uint              `json:"id"`
	CustomerID        uint              `json:"customer_id"`
	EmployeeID        *uint             `json:"employee_id,omitempty"`
	StoreID           *uint             `json:"store_id,omitempty"`
	DeliveryAddressID *uint             `json:"delivery_address_id,omitempty"`
	OrderStatus       string            `json:"order_status"`
	OrderDate         time.Time         `json:"order_date"`
	Total             float64           `json:"total"`
	OrderProducts     []OrderProductDTO `json:"order_products,omitempty"`
}

type OrderProductDTO struct {
	ID             uint                      `json:"id"`
	OrderID        uint                      `json:"order_id"`
	ProductSizeID  uint                      `json:"product_size_id"`
	Quantity       int                       `json:"quantity"`
	Price          float64                   `json:"price"`
	SubOrderStatus string                    `json:"sub_order_status"`
	Additives      []OrderProductAdditiveDTO `json:"additives,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

type OrderProductAdditiveDTO struct {
	ID             uint      `json:"id"`
	OrderProductID uint      `json:"order_product_id"`
	AdditiveID     uint      `json:"additive_id"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
