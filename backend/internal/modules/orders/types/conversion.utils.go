package types

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ConvertCreateOrderDTOToOrder(createOrderDTO *CreateOrderDTO, productPrices map[uint]float64, additivePrices map[uint]float64) (data.Order, float64) {
	var total float64
	var orderProducts []data.OrderProduct

	for _, productDTO := range createOrderDTO.OrderItems {
		var additives []data.OrderProductAdditive
		var productTotal float64

		if price, exists := productPrices[productDTO.ProductSizeID]; exists {
			productTotal = price * float64(productDTO.Quantity)
		}

		for _, additiveID := range productDTO.AdditivesIDs {
			if price, exists := additivePrices[additiveID]; exists {
				additives = append(additives, data.OrderProductAdditive{
					AdditiveID: additiveID,
					Price:      price,
				})
				productTotal += price
			}
		}

		orderProducts = append(orderProducts, data.OrderProduct{
			ProductSizeID: productDTO.ProductSizeID,
			Quantity:      productDTO.Quantity,
			Price:         productTotal,
			Additives:     additives,
		})

		total += productTotal
	}

	order := data.Order{
		CustomerID:        createOrderDTO.CustomerID,
		CustomerName:      createOrderDTO.CustomerName,
		EmployeeID:        createOrderDTO.EmployeeID,
		StoreID:           createOrderDTO.StoreID,
		DeliveryAddressID: createOrderDTO.DeliveryAddressID,
		OrderStatus:       data.OrderStatusPending,
		OrderProducts:     orderProducts,
	}

	return order, total
}

func ConvertOrderToDTO(order *data.Order) OrderDTO {
	orderDTO := OrderDTO{
		ID:                order.ID,
		CustomerID:        order.CustomerID,
		CustomerName:      &order.CustomerName,
		EmployeeID:        order.EmployeeID,
		StoreID:           order.StoreID,
		DeliveryAddressID: order.DeliveryAddressID,
		OrderStatus:       order.OrderStatus,
		CreatedAt:         order.CreatedAt,
		Total:             order.Total,
	}

	for _, product := range order.OrderProducts {
		orderDTO.OrderProducts = append(orderDTO.OrderProducts, ConvertOrderProductToDTO(&product))
	}

	return orderDTO
}

func ConvertOrderProductToDTO(product *data.OrderProduct) OrderProductDTO {
	productDTO := OrderProductDTO{
		ID:            product.ID,
		OrderID:       product.OrderID,
		ProductSizeID: product.ProductSizeID,
		Quantity:      product.Quantity,
		Price:         product.Price,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}

	for _, additive := range product.Additives {
		productDTO.Additives = append(productDTO.Additives, ConvertOrderProductAdditiveToDTO(&additive))
	}

	return productDTO
}

func ConvertOrderProductAdditiveToDTO(additive *data.OrderProductAdditive) OrderProductAdditiveDTO {
	return OrderProductAdditiveDTO{
		ID:             additive.ID,
		OrderProductID: additive.OrderProductID,
		AdditiveID:     additive.AdditiveID,
		Price:          additive.Price,
		CreatedAt:      additive.CreatedAt,
		UpdatedAt:      additive.UpdatedAt,
	}
}

func ConvertDTOToOrder(orderDTO *OrderDTO) data.Order {
	order := data.Order{
		CustomerID:        orderDTO.CustomerID,
		EmployeeID:        orderDTO.EmployeeID,
		StoreID:           orderDTO.StoreID,
		DeliveryAddressID: orderDTO.DeliveryAddressID,
		OrderStatus:       orderDTO.OrderStatus,
		Total:             orderDTO.Total,
	}

	for _, productDTO := range orderDTO.OrderProducts {
		order.OrderProducts = append(order.OrderProducts, ConvertDTOToOrderProduct(&productDTO))
	}

	return order
}

func ConvertDTOToOrderProduct(productDTO *OrderProductDTO) data.OrderProduct {
	product := data.OrderProduct{
		OrderID:       productDTO.OrderID,
		ProductSizeID: productDTO.ProductSizeID,
		Quantity:      productDTO.Quantity,
		Price:         productDTO.Price,
	}

	for _, additiveDTO := range productDTO.Additives {
		product.Additives = append(product.Additives, ConvertDTOToOrderProductAdditive(&additiveDTO))
	}

	return product
}

func ConvertDTOToOrderProductAdditive(additiveDTO *OrderProductAdditiveDTO) data.OrderProductAdditive {
	return data.OrderProductAdditive{
		OrderProductID: additiveDTO.OrderProductID,
		AdditiveID:     additiveDTO.AdditiveID,
		Price:          additiveDTO.Price,
	}
}

func ArrayToCSV(arr []uint) string {
	var result string
	for i, val := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", val)
	}
	return result
}
