package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type OrderEvent struct {
	OrderID   uint             `json:"order_id"`
	StoreID   *uint            `json:"store_id"`
	Status    data.OrderStatus `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Items     []SubOrderEvent  `json:"items"`
}

type SubOrderEvent struct {
	SubOrderID uint             `json:"sub_order_id"`
	Status     data.OrderStatus `json:"status"`
}
