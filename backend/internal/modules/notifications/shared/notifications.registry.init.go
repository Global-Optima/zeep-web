package shared

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
)

func InitNotificationRegistry() {
	RegisterNotification(
		data.CENTRAL_CATALOG_UPDATE,
		func() details.NotificationDetails {
			return &details.CentralCatalogUpdateDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			catalogDetails, ok := baseDetails.(*details.CentralCatalogUpdateDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for CENTRAL_CATALOG_UPDATE")
			}
			return details.BuildCentralCatalogUpdateMessage(catalogDetails)
		},
	)

	RegisterNotification(
		data.NEW_ORDER,
		func() details.NotificationDetails {
			return &details.NewOrderNotificationDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			orderDetails, ok := baseDetails.(*details.NewOrderNotificationDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for NEW_ORDER")
			}
			return details.BuildNewOrderMessage(orderDetails)
		},
	)

	RegisterNotification(
		data.NEW_STOCK_REQUEST,
		func() details.NotificationDetails {
			return &details.NewStockRequestDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			stockRequestDetails, ok := baseDetails.(*details.NewStockRequestDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for NEW_STOCK_REQUEST")
			}
			return details.BuildNewStockRequestMessage(stockRequestDetails)
		},
	)

	RegisterNotification(
		data.WAREHOUSE_OUT_OF_STOCK,
		func() details.NotificationDetails {
			return &details.OutOfStockDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			outOfStockDetails, ok := baseDetails.(*details.OutOfStockDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for OUT_OF_STOCK")
			}
			return details.BuildOutOfStockMessage(outOfStockDetails)
		},
	)

	RegisterNotification(
		data.PRICE_CHANGE,
		func() details.NotificationDetails {
			return &details.PriceChangeNotificationDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			priceChangeDetails, ok := baseDetails.(*details.PriceChangeNotificationDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for PRICE_CHANGE")
			}
			return details.BuildPriceChangeMessage(priceChangeDetails)
		},
	)

	RegisterNotification(
		data.WAREHOUSE_STOCK_EXPIRATION,
		func() details.NotificationDetails {
			return &details.StockExpirationDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			stockExpirationDetails, ok := baseDetails.(*details.StockExpirationDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for STOCK_EXPIRATION")
			}
			return details.BuildStockExpirationMessage(stockExpirationDetails)
		},
	)

	RegisterNotification(
		data.STOCK_REQUEST_STATUS_UPDATED,
		func() details.NotificationDetails {
			return &details.StockRequestStatusUpdatedDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			requestDetails, ok := baseDetails.(*details.StockRequestStatusUpdatedDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for STOCK_REQUEST_STATUS_UPDATED")
			}
			return details.BuildStockRequestStatusUpdatedMessage(requestDetails)
		},
	)

	RegisterNotification(
		data.STORE_WAREHOUSE_RUN_OUT,
		func() details.NotificationDetails {
			return &details.StoreWarehouseRunOutDetails{}
		},
		func(baseDetails details.NotificationDetails) (localization.LocalizedMessage, error) {
			warehouseDetails, ok := baseDetails.(*details.StoreWarehouseRunOutDetails)
			if !ok {
				return localization.LocalizedMessage{}, fmt.Errorf("invalid details type for STORE_WAREHOUSE_RUN_OUT")
			}
			return details.BuildStoreWarehouseRunOutMessage(warehouseDetails)
		},
	)
}
