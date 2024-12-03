package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type OrderEvent struct {
	ID                uint            `json:"id"`
	CustomerName      string          `json:"customer_name"`
	ETA               time.Time       `json:"eta"`
	OrderType         string          `json:"order_type"` // "Доставка" or "Кафе"
	StoreID           uint            `json:"store_id"`
	SubOrdersQuantity int             `json:"sub_orders_quantity"`
	Items             []SubOrderEvent `json:"items"`
	Status            string          `json:"status"`
	Timestamp         time.Time       `json:"timestamp"`
}

type SubOrderEvent struct {
	SubOrderID    uint             `json:"sub_order_id"`
	ProductSizeID uint             `json:"product_size_id"`
	ProductName   string           `json:"product_name"`
	Additives     []AdditiveDetail `json:"additives"`
	ETA           time.Time        `json:"eta"`
	Status        data.OrderStatus `json:"status"`
}

type AdditiveDetail struct {
	AdditiveID uint    `json:"additive_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
}
