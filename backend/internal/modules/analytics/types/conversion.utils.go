package types

func ToSummaryDTO(totalSales float64, totalOrders, totalProductsSold, totalAdditivesSold int, previousMonthSales float64, previousMonthOrders int) SummaryDTO {
	salesComparison := 0.0
	if previousMonthSales > 0 {
		salesComparison = ((totalSales - previousMonthSales) / previousMonthSales) * 100
	}

	ordersComparison := 0.0
	if previousMonthOrders > 0 {
		ordersComparison = float64(totalOrders-previousMonthOrders) / float64(previousMonthOrders) * 100
	}

	return SummaryDTO{
		TotalSales:          totalSales,
		TotalOrders:         totalOrders,
		TotalProductsSold:   totalProductsSold,
		TotalAdditivesSold:  totalAdditivesSold,
		PreviousMonthSales:  previousMonthSales,
		PreviousMonthOrders: previousMonthOrders,
		SalesComparison:     salesComparison,
		OrdersComparison:    ordersComparison,
	}
}

func ToMonthlySalesDTO(month string, year int, sales float64, orders int) MonthlySalesDTO {
	return MonthlySalesDTO{
		Month:  month,
		Year:   year,
		Sales:  sales,
		Orders: orders,
	}
}

func ToPopularProductDTO(productName string, totalSold int, revenue float64) PopularProductDTO {
	return PopularProductDTO{
		ProductName: productName,
		TotalSold:   totalSold,
		Revenue:     revenue,
	}
}

func ToProductSoldDTO(productName string, totalSold int, totalRevenue float64, additives []AdditiveSoldDTO) ProductSoldDTO {
	return ProductSoldDTO{
		ProductName:   productName,
		TotalSold:     totalSold,
		TotalRevenue:  totalRevenue,
		AdditivesSold: additives,
	}
}

func ToAdditiveSoldDTO(additiveName string, quantity int, revenue float64) AdditiveSoldDTO {
	return AdditiveSoldDTO{
		AdditiveName: additiveName,
		Quantity:     quantity,
		Revenue:      revenue,
	}
}
