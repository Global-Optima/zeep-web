package types

import (
	"time"
)

type OrderEvent struct {
	OrderID   uint            `json:"order_id"`
	StoreID   *uint           `json:"store_id"`
	Status    string          `json:"status"`
	Timestamp time.Time       `json:"timestamp"`
	Items     []SubOrderEvent `json:"items"`
}

type SubOrderEvent struct {
	SubOrderID uint   `json:"sub_order_id"`
	Status     string `json:"status"` // e.g., 'PENDING', 'COMPLETED'
}
