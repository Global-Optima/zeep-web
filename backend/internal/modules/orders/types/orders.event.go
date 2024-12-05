package types

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

type OrderEvent struct {
	ID                uint            `json:"id"`
	CustomerName      string          `json:"customerName"`
	ETA               time.Time       `json:"eta"`
	OrderType         string          `json:"orderType"` // "Доставка" or "Кафе"
	StoreID           uint            `json:"storeId"`
	SubOrdersQuantity int             `json:"subOrdersQuantity"`
	Items             []SubOrderEvent `json:"items"`
	Status            string          `json:"status"`
	Timestamp         time.Time       `json:"timestamp"`
}

type SubOrderEvent struct {
	ID            uint             `json:"id"`
	SubOrderID    uint             `json:"subOrderId"`
	ProductSizeID uint             `json:"productSizeId"`
	ProductName   string           `json:"productName"`
	Additives     []AdditiveDetail `json:"additives"`
	ETA           time.Time        `json:"eta"`
	Status        data.OrderStatus `json:"status"`
}

type AdditiveDetail struct {
	AdditiveID uint    `json:"additiveId"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
}
