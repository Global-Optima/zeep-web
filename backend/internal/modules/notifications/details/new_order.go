package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

// NewOrderNotificationDetails specific to new order notifications
type NewOrderNotificationDetails struct {
	BaseNotificationDetails
	CustomerName string `json:"customerName"`
	OrderID      uint   `json:"orderId"`
}

func (n *NewOrderNotificationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NewOrderNotificationDetails) GetBaseDetails() *BaseNotificationDetails {
	return &n.BaseNotificationDetails
}

func BuildNewOrderDetails(facilityID uint, facilityName, customerName string, orderID uint) (*NewOrderNotificationDetails, error) {
	if facilityID == 0 || orderID == 0 || facilityName == "" || customerName == "" {
		return nil, fmt.Errorf("invalid input: facilityID, orderID, facilityName, and customerName are required")
	}

	return &NewOrderNotificationDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		CustomerName: customerName,
		OrderID:      orderID,
	}, nil
}

func BuildNewOrderMessage(details *NewOrderNotificationDetails) (localization.LocalizedMessages, error) {
	if details == nil {
		return localization.LocalizedMessages{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.NEW_ORDER.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName": details.FacilityName,
		"ID":           details.ID,
		"CustomerName": details.CustomerName,
		"OrderID":      details.OrderID,
	})
	if err != nil {
		return localization.LocalizedMessages{}, fmt.Errorf("failed to build NewOrder message: %w", err)
	}

	return *messages, nil
}
