package types

import (
	"time"
)

type SummaryDTO struct {
	TotalSales          float64 `json:"totalSales"`
	TotalOrders         int     `json:"totalOrders"`
	TotalProductsSold   int     `json:"totalProductsSold"`
	TotalAdditivesSold  int     `json:"totalAdditivesSold"`
	PreviousMonthSales  float64 `json:"previousMonthSales"`
	PreviousMonthOrders int     `json:"previousMonthOrders"`
	SalesComparison     float64 `json:"salesComparison"`
	OrdersComparison    float64 `json:"ordersComparison"`
}

type MonthlySalesDTO struct {
	Month  string  `json:"month"`
	Year   int     `json:"year"`
	Sales  float64 `json:"sales"`
	Orders int     `json:"orders"`
}

type PopularProductDTO struct {
	ProductName string  `json:"productName"`
	TotalSold   int     `json:"totalSold"`
	Revenue     float64 `json:"revenue"`
}

type ProductSoldDTO struct {
	ProductName   string            `json:"productName"`
	TotalSold     int               `json:"totalSold"`
	TotalRevenue  float64           `json:"totalRevenue"`
	AdditivesSold []AdditiveSoldDTO `json:"additivesSold"`
}

type AdditiveSoldDTO struct {
	AdditiveName string  `json:"additiveName"`
	Quantity     int     `json:"quantity"`
	Revenue      float64 `json:"revenue"`
}

type AnalyticsFilterQuery struct {
	StartDate time.Time `form:"startDate" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time `form:"endDate" binding:"required" time_format:"2006-01-02"`
	StoreID   uint      `form:"storeId" binding:"required" time_format:"2006-01-02"`
}
