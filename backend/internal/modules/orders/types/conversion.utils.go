package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
)

func ConvertCreateOrderDTOToOrder(createOrderDTO *CreateOrderDTO, productPrices map[uint]float64, additivePrices map[uint]float64) (data.Order, float64) {
	var total float64
	var suborders []data.Suborder

	for _, productDTO := range createOrderDTO.Suborders {
		for i := 0; i < productDTO.Quantity; i++ {
			var additives []data.SuborderAdditive
			suborderTotal := productPrices[productDTO.ProductSizeID]

			for _, additiveID := range productDTO.AdditivesIDs {
				if price, exists := additivePrices[additiveID]; exists {
					additives = append(additives, data.SuborderAdditive{
						AdditiveID: additiveID,
						Price:      price,
					})
					suborderTotal += price
				}
			}

			suborders = append(suborders, data.Suborder{
				ProductSizeID: productDTO.ProductSizeID,
				Price:         suborderTotal,
				Status:        data.SubOrderStatusPreparing,
				Additives:     additives,
			})

			total += suborderTotal
		}
	}

	return data.Order{
		CustomerID:        createOrderDTO.CustomerID,
		CustomerName:      createOrderDTO.CustomerName,
		EmployeeID:        createOrderDTO.EmployeeID,
		StoreID:           createOrderDTO.StoreID,
		DeliveryAddressID: createOrderDTO.DeliveryAddressID,
		Status:            data.OrderStatusPreparing,
		Suborders:         suborders,
	}, total
}

func ConvertOrderToDTO(order *data.Order) OrderDTO {
	orderDTO := OrderDTO{
		ID:                order.ID,
		CustomerID:        order.CustomerID,
		CustomerName:      &order.CustomerName,
		EmployeeID:        order.EmployeeID,
		StoreID:           order.StoreID,
		DeliveryAddressID: order.DeliveryAddressID,
		Status:            order.Status,
		CreatedAt:         order.CreatedAt,
		Total:             order.Total,
		SubordersQuantity: len(order.Suborders),
		Suborders:         []SuborderDTO{},
	}

	for _, suborder := range order.Suborders {
		orderDTO.Suborders = append(orderDTO.Suborders, ConvertSuborderToDTO(&suborder))
	}

	return orderDTO
}

func ConvertSuborderToDTO(suborder *data.Suborder) SuborderDTO {
	suborderDTO := SuborderDTO{
		ID:      suborder.ID,
		OrderID: suborder.OrderID,
		ProductSize: OrderProductSizeDTO{
			ID:          suborder.ProductSize.ID,
			SizeName:    suborder.ProductSize.Name,
			ProductName: suborder.ProductSize.Product.Name,
			Size:        suborder.ProductSize.Size,
			Unit:        unitTypes.ToUnitResponse(suborder.ProductSize.Unit),
		},
		Price:     suborder.Price,
		Status:    suborder.Status,
		CreatedAt: suborder.CreatedAt,
		UpdatedAt: suborder.UpdatedAt,
		Additives: []SuborderAdditiveDTO{},
	}

	for _, additive := range suborder.Additives {
		suborderDTO.Additives = append(suborderDTO.Additives, ConvertSuborderAdditiveToDTO(&additive))
	}

	return suborderDTO
}

func ConvertSuborderAdditiveToDTO(additive *data.SuborderAdditive) SuborderAdditiveDTO {
	return SuborderAdditiveDTO{
		ID:         additive.ID,
		SuborderID: additive.SuborderID,
		Additive: OrderAdditiveDTO{
			ID:          additive.Additive.ID,
			Name:        additive.Additive.Name,
			Description: additive.Additive.Description,
			Size:        additive.Additive.Size,
		},
		Price:     additive.Price,
		CreatedAt: additive.CreatedAt,
		UpdatedAt: additive.UpdatedAt,
	}
}

func MapToOrderDetailsDTO(order *data.Order) *OrderDetailsDTO {
	if order == nil {
		return nil
	}

	suborders := make([]SuborderDetailsDTO, len(order.Suborders))
	for i, sub := range order.Suborders {
		additives := make([]OrderAdditiveDetailsDTO, len(sub.Additives))
		for j, additive := range sub.Additives {
			additives[j] = OrderAdditiveDetailsDTO{
				ID:          additive.Additive.ID,
				Name:        additive.Additive.Name,
				Description: additive.Additive.Description,
				BasePrice:   additive.Additive.BasePrice,
			}
		}

		suborders[i] = SuborderDetailsDTO{
			ID:     sub.ID,
			Price:  sub.Price,
			Status: string(sub.Status),
			ProductSize: OrderProductSizeDetailsDTO{
				ID:        sub.ProductSize.ID,
				Name:      sub.ProductSize.Name,
				Unit:      unitTypes.ToUnitResponse(sub.ProductSize.Unit),
				BasePrice: sub.ProductSize.BasePrice,
				Product: OrderProductDetailsDTO{
					ID:          sub.ProductSize.Product.ID,
					Name:        sub.ProductSize.Product.Name,
					Description: sub.ProductSize.Product.Description,
					ImageURL:    sub.ProductSize.Product.ImageURL,
				},
			},
			Additives: additives,
		}
	}

	var deliveryAddress *OrderDeliveryAddressDTO
	if order.DeliveryAddressID != nil {
		deliveryAddress = &OrderDeliveryAddressDTO{
			ID:        order.DeliveryAddress.ID,
			Address:   order.DeliveryAddress.Address,
			Longitude: order.DeliveryAddress.Longitude,
			Latitude:  order.DeliveryAddress.Latitude,
		}
	}

	var customerName *string
	if order.CustomerName != "" {
		customerName = &order.CustomerName
	}

	return &OrderDetailsDTO{
		ID:              order.ID,
		CustomerName:    customerName, // Optional
		Status:          string(order.Status),
		Total:           order.Total,
		Suborders:       suborders,
		DeliveryAddress: deliveryAddress, // Optional
	}
}

// for export
func ToOrderExportDTO(order *data.Order, storeName string) OrderExportDTO {
	suborders := make([]SuborderDTO, len(order.Suborders))
	for i, suborder := range order.Suborders {
		suborders[i] = ConvertSuborderToDTO(&suborder)
	}

	return OrderExportDTO{
		ID:              order.ID,
		CustomerName:    order.CustomerName,
		Status:          string(order.Status),
		Total:           order.Total,
		CreatedAt:       order.CreatedAt,
		StoreName:       storeName,
		Suborders:       suborders,
		DeliveryAddress: ToOrderDeliveryAddressDTO(&order.DeliveryAddress),
	}
}

func ToOrderAdditivesDTO(additives []data.SuborderAdditive) []SuborderAdditiveDTO {
	dtos := make([]SuborderAdditiveDTO, len(additives))
	for i, additive := range additives {
		dtos[i] = ConvertSuborderAdditiveToDTO(&additive)
	}
	return dtos
}

func ToOrderDeliveryAddressDTO(address *data.CustomerAddress) *OrderDeliveryAddressDTO {
	if address == nil {
		return nil
	}
	return &OrderDeliveryAddressDTO{
		ID:        address.ID,
		Address:   address.Address,
		Longitude: address.Longitude,
		Latitude:  address.Latitude,
	}
}
