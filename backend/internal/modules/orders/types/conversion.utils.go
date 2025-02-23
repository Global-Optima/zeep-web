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
			var storeAdditives []data.SuborderAdditive
			suborderTotal := productPrices[productDTO.StoreProductSizeID]

			for _, additiveID := range productDTO.StoreAdditivesIDs {
				if price, exists := additivePrices[additiveID]; exists {
					storeAdditives = append(storeAdditives, data.SuborderAdditive{
						StoreAdditiveID: additiveID,
						Price:           price,
					})
					suborderTotal += price
				}
			}

			suborders = append(suborders, data.Suborder{
				StoreProductSizeID: productDTO.StoreProductSizeID,
				Price:              suborderTotal,
				Status:             data.SubOrderStatusPending,
				SuborderAdditives:  storeAdditives,
			})

			total += suborderTotal
		}
	}

	return data.Order{
		CustomerID:        createOrderDTO.CustomerID,
		CustomerName:      createOrderDTO.CustomerName,
		StoreEmployeeID:   createOrderDTO.StoreEmployeeID,
		StoreID:           createOrderDTO.StoreID,
		DeliveryAddressID: createOrderDTO.DeliveryAddressID,
		Status:            data.OrderStatusPending,
		Suborders:         suborders,
	}, total
}

func ConvertOrderToDTO(order *data.Order) OrderDTO {
	orderDTO := OrderDTO{
		ID:                order.ID,
		CustomerID:        order.CustomerID,
		CustomerName:      &order.CustomerName,
		StoreEmployeeID:   order.StoreEmployeeID,
		StoreID:           order.StoreID,
		DeliveryAddressID: order.DeliveryAddressID,
		Status:            order.Status,
		CreatedAt:         order.CreatedAt,
		Total:             order.Total,
		SubordersQuantity: len(order.Suborders),
		Suborders:         []SuborderDTO{},
		DisplayNumber:     order.DisplayNumber,
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
		ProductSize: OrderStoreProductSizeDTO{
			ID:          suborder.StoreProductSize.ProductSize.ID,
			SizeName:    suborder.StoreProductSize.ProductSize.Name,
			ProductName: suborder.StoreProductSize.ProductSize.Product.Name,
			Size:        suborder.StoreProductSize.ProductSize.Size,
			Unit:        unitTypes.ToUnitResponse(suborder.StoreProductSize.ProductSize.Unit),
		},
		Price:     suborder.Price,
		Status:    suborder.Status,
		CreatedAt: suborder.CreatedAt,
		UpdatedAt: suborder.UpdatedAt,
		Additives: []SuborderStoreAdditiveDTO{},
	}

	for _, additive := range suborder.SuborderAdditives {
		suborderDTO.Additives = append(suborderDTO.Additives, ConvertSuborderAdditiveToDTO(&additive))
	}

	return suborderDTO
}

func ConvertSuborderAdditiveToDTO(suborderAdditive *data.SuborderAdditive) SuborderStoreAdditiveDTO {
	return SuborderStoreAdditiveDTO{
		ID:         suborderAdditive.ID,
		SuborderID: suborderAdditive.SuborderID,
		Additive: OrderStoreAdditiveDTO{
			ID:          suborderAdditive.StoreAdditive.Additive.ID,
			Name:        suborderAdditive.StoreAdditive.Additive.Name,
			Description: suborderAdditive.StoreAdditive.Additive.Description,
			Size:        suborderAdditive.StoreAdditive.Additive.Size,
		},
		Price:     suborderAdditive.Price,
		CreatedAt: suborderAdditive.CreatedAt,
		UpdatedAt: suborderAdditive.UpdatedAt,
	}
}

func MapToOrderDetailsDTO(order *data.Order) *OrderDetailsDTO {
	if order == nil {
		return nil
	}

	suborders := make([]SuborderDetailsDTO, len(order.Suborders))
	for i, sub := range order.Suborders {
		storeAdditives := make([]OrderAdditiveDetailsDTO, len(sub.SuborderAdditives))
		for j, storeAdditive := range sub.SuborderAdditives {
			storeAdditives[j] = OrderAdditiveDetailsDTO{
				ID:          storeAdditive.StoreAdditive.Additive.ID,
				Name:        storeAdditive.StoreAdditive.Additive.Name,
				Description: storeAdditive.StoreAdditive.Additive.Description,
				BasePrice:   storeAdditive.StoreAdditive.Additive.BasePrice,
			}
		}

		suborders[i] = SuborderDetailsDTO{
			ID:     sub.ID,
			Price:  sub.Price,
			Status: string(sub.Status),
			StoreProductSize: OrderProductSizeDetailsDTO{
				ID:         sub.StoreProductSize.ID,
				Name:       sub.StoreProductSize.ProductSize.Name,
				Unit:       unitTypes.ToUnitResponse(sub.StoreProductSize.ProductSize.Unit),
				StorePrice: sub.StoreProductSize.ProductSize.BasePrice,
				Product: OrderProductDetailsDTO{
					ID:          sub.StoreProductSize.ProductSize.Product.ID,
					Name:        sub.StoreProductSize.ProductSize.Product.Name,
					Description: sub.StoreProductSize.ProductSize.Product.Description,
					ImageURL:    sub.StoreProductSize.ProductSize.Product.ImageURL.GetURL(),
				},
			},
			StoreAdditives: storeAdditives,
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

func ToOrderAdditivesDTO(additives []data.SuborderAdditive) []SuborderStoreAdditiveDTO {
	dtos := make([]SuborderStoreAdditiveDTO, len(additives))
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
